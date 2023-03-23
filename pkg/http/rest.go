package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/cors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const timeout = 60

// RestConfig holds all necessary configuration to run module rest.
type RestConfig struct {
	CorsAllowedOrigins  []string `split_words:"true" required:"false" default:"*"`
	CorsAllowedHeaders  []string `split_words:"true" required:"false" default:"Accept,Authorization,Content-Type"`
	CorsAllowedMehthods []string `split_words:"true" required:"false" default:"GET,POST,HEAD"`
}

// ProvideRestConfig process the configuration needed to run a rest module.
func ProvideRestConfig(l *zap.Logger) (*RestConfig, error) {
	var config RestConfig
	if err := envconfig.Process("rest", &config); err != nil {
		l.Error(ErrEnvConfig.Error(), zap.Error(err))

		return nil, ErrEnvConfig
	}

	return &config, nil
}

// ProvideRouter provides a new instance of an HTTP mux from go-chi/chi package.
func ProvideRouter(l *zap.Logger, config *RestConfig) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.New(cors.Options{
		AllowedOrigins: config.CorsAllowedOrigins,
		AllowedHeaders: config.CorsAllowedHeaders,
		AllowedMethods: config.CorsAllowedMehthods,
	}).Handler)
	r.Use(middleware.AllowContentEncoding("application/json"))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(logMiddleware(l))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	return r
}

// ServePlainHTTP starts a simple HTTP server and configures a hook lifecycle for working
// on top of uber/fx package.
func ServePlainHTTP(lc fx.Lifecycle, handler *chi.Mux, logger *zap.Logger, listener net.Listener) error {
	server := &http.Server{
		ReadHeaderTimeout: timeout * time.Second,
		Handler:           handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("starting HTTP server...", zap.String("address", listener.Addr().String()))

			// NOTE: this error is not checked. Fix in the future. More details in here:
			// https://github.com/uber-go/fx/issues/600
			//nolint:errcheck
			go server.Serve(listener)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping HTTP server...")

			err := server.Shutdown(ctx)
			if err != nil {
				logger.Error(ErrHTTPCloseConn.Error(), zap.Error(err))

				return ErrHTTPCloseConn
			}

			return nil
		},
	})

	return nil
}

// logMiddleware is a middleware for logging useful information
// of each request, which makes use of uber/zap package.
func logMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&loggerRequest{log: logger})
}

// loggerRequest implements middleware.LogFormatter interface.
type loggerRequest struct {
	log *zap.Logger
}

// NewLogEntry creates a new LogEntry for the request.
func (l *loggerRequest) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	ctx := r.Context()

	entryLogger := l.log.With(
		zap.String("request_id", middleware.GetReqID(ctx)),
		zap.String("http_uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
	)
	entryLogger.Info("http request started",
		zap.String("http_scheme", scheme),
		zap.String("http_protocol", r.Proto),
		zap.String("http_method", r.Method),
		zap.String("remote_address", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
	)

	return &loggerRequestEntry{log: entryLogger}
}

// loggerRequestEntry implements middleware.LogEntry interface.
type loggerRequestEntry struct {
	log *zap.Logger
}

// Write is called whenever a request is completed and prints
// useful information about the result.
func (l *loggerRequestEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.log.Check(logLevel(status), "http request completed").Write(
		zap.Int("resp_status", status),
		zap.Int("resp_bytes_length", bytes),
		zap.String("elapsed_time", elapsed.String()),
	)
}

// Panic is called whenever the application panics, with the provided message.
func (l *loggerRequestEntry) Panic(v interface{}, stack []byte) {
	l.log.Panic("fatal error", zap.String("panic_message", fmt.Sprintf("%+v", v)))
}

// logLevel is an utility function to map HTTP log levels
// to standard levels of uber/zap Logger.
// https://developer.mozilla.org/pt-BR/docs/Web/HTTP/Status
func logLevel(status int) zapcore.Level {
	switch {
	case status >= 100 && status < 500:
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}

// WithAPM returns a handler instrumented for observability.
func WithAPM(cli *newrelic.Application, pattern string, handler http.HandlerFunc) http.HandlerFunc {
	if cli != nil {
		_, handler = newrelic.WrapHandleFunc(cli, pattern, handler)
	}

	return handler
}

// RESTModule wrapper for uber/fx.
//
//nolint:gochecknoglobals
var RESTModule = fx.Options(
	httpModule,
	fx.Provide(ProvideRestConfig),
	fx.Provide(ProvideRouter),
	fx.Invoke(ServePlainHTTP),
)
