// Package log implements go.uber.org/zap.Logger and allows
// customization to work properly with Cloud Logging.
// This package will initialize a global instance of the Logger.
package log

import (
	"context"
	"encoding/json"

	"github.com/blendle/zapdriver"
	"github.com/gogo/protobuf/version"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// loggerConfig holds all configuration from this package.
type loggerConfig struct {
	GCP struct {
		EnableStackdriver    bool   `split_words:"true" required:"false" default:"true"`
		EnableErrorReporting bool   `split_words:"true" required:"false" default:"true"`
		ServiceName          string `split_words:"true" required:"false"`
	}
	DebugLevel      bool `split_words:"true" required:"false" default:"false"`
	DevelopmentMode bool `split_words:"true" required:"false" default:"false"`
}

// ProvideLogger returns an instance of *zap.Logger.
//
// It panics by design if some error occurs during Logger setup.
// This is the expected behavior, since every other service
// in our ecosystem will rely on this to do logging.
//
// By default, Logger will be wrapped into the zapdriver package.
// This feature can be disabled through the use of
// `LOG_GCP_ENABLE_STACKDRIVER` environment variable.
//
// It is also possible to customize the driver for working
// properly with Cloud Logging (former Stackdriver).
// The following environment variables can be set:
// - `LOG_GCP_ENABLE_ERROR_REPORTING` (default true)
// - `LOG_GCP_SERVICE_NAME` (default empty).
func ProvideLogger() *zap.Logger {
	var lc loggerConfig
	if err := envconfig.Process("log", &lc); err != nil {
		panic(err)
	}

	var opts []zap.Option
	if !lc.DevelopmentMode && lc.GCP.EnableStackdriver {
		// Error reporting can be enabled for working with GCP Error Reporting.
		// More info: https://cloud.google.com/error-reporting/docs
		opts = append(opts, zapdriver.WrapCore(
			zapdriver.ReportAllErrors(lc.GCP.EnableErrorReporting),
			zapdriver.ServiceName(lc.GCP.ServiceName),
		))
	}

	logger, err := config(lc).Build(opts...)
	if err != nil {
		panic(err)
	}

	return logger
}

// RegisterLoggerOnStopHook is used to register a hook to be executed
// whenever the application stops. This will flush the logs.
func RegisterLoggerOnStopHook(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("flushing logger...")
			_ = logger.Sync()

			return nil
		},
	})
}

// config setups default production configuration from
// go.uber.org/zap.Logger.
//
// It allows configuring the wanted debug level
// through environment variable `LOG_DEBUG_LEVEL`.
//
// Default value does not print debug logs.
func config(lc loggerConfig) zap.Config {
	if lc.DevelopmentMode {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		return cfg
	}

	cfg := zap.NewProductionConfig()
	if lc.DebugLevel {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	if lc.GCP.EnableStackdriver {
		cfg.EncoderConfig = zapdriver.NewProductionEncoderConfig()
	}

	var initialFields map[string]interface{}

	versionBytes, err := json.Marshal(version.Get())
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(versionBytes, &initialFields)
	if err != nil {
		panic(err)
	}

	cfg.InitialFields = initialFields

	return cfg
}

// Module wrapper for uber/fx.
//
//nolint:gochecknoglobals
var Module = fx.Options(
	fx.Provide(ProvideLogger),
	fx.Invoke(RegisterLoggerOnStopHook),
	fx.WithLogger(
		func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		},
	),
)
