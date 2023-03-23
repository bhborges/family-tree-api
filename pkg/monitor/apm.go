package monitor

import (
	"context"

	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// APMConfig holds all necessary configuration to run an New Relic client.
type APMConfig struct {
	AppName    string `split_words:"true" required:"false"`
	LicenseKey string `split_words:"true" required:"false"`
}

// ProvideAPMConfig process the configuration needed to access GSC.
func ProvideAPMConfig(log *zap.Logger) (*APMConfig, error) {
	var config APMConfig

	err := envconfig.Process("monitor", &config)
	if err != nil {
		return nil, ErrEnvConfig
	}

	return &config, nil
}

// ProvideNewRelicApp create a app instance for new relic with the given options.
func ProvideNewRelicApp(conf *APMConfig, log *zap.Logger) (*newrelic.Application, error) {
	if (conf.AppName == "") || (conf.LicenseKey == "") {
		log.Warn("starting application without new relic observability...")

		return nil, nil
	}

	newrelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(conf.AppName),
		newrelic.ConfigLicense(conf.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return nil, ErrProvideApp
	}

	return newrelicApp, nil
}

// RegisterAPMOnStopHook is used to register a hook to be executed
// whenever the application stops. This will flush the logs.
func RegisterAPMOnStopHook(lc fx.Lifecycle, newrelicApp *newrelic.Application) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if newrelicApp != nil {
				newrelicApp.Shutdown(0)
			}

			return nil
		},
	})
}

// APMModule wrapper for uber/fx.
func APMModule() fx.Option {
	return fx.Options(
		fx.Provide(ProvideAPMConfig),
		fx.Provide(ProvideNewRelicApp),
		fx.Invoke(RegisterAPMOnStopHook),
	)
}
