// Package http uses go-chi/chi for routing when building an HTTP server.
// and provides a simple wrapper to start an HTTP server using uber/fx.
package http

import (
	"fmt"
	"net"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ListenerConfig holds all necessary configuration to
// run an HTTP server.
type ListenerConfig struct {
	ServerPort string `split_words:"true" required:"false" default:"5001"`
}

// ProvideListenerConfig process the configuration needed
// to run an HTTP server.
func ProvideListenerConfig(l *zap.Logger) (*ListenerConfig, error) {
	var config ListenerConfig
	if err := envconfig.Process("http", &config); err != nil {
		l.Error(ErrEnvConfig.Error(), zap.Error(err))

		return nil, ErrEnvConfig
	}

	return &config, nil
}

// ProvideTCPListener announces the TCP connection and returns
// an instance of a net.Listener object.
func ProvideTCPListener(l *zap.Logger, config *ListenerConfig) (net.Listener, error) {
	addr := fmt.Sprintf(":%s", config.ServerPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error(ErrTCPListening.Error(), zap.Error(err))

		return nil, ErrTCPListening
	}

	l.Debug("TCP connection successfully announced", zap.String("address", addr))

	return listener, nil
}

// httpModule wrapper for uber/fx.
//
//nolint:gochecknoglobals
var httpModule = fx.Options(
	fx.Provide(ProvideListenerConfig),
	fx.Provide(ProvideTCPListener),
)
