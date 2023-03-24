// Package family tree holds all modules related to it.
package familytree

import (
	"github.com/bhborges/family-tree-api/internal/adapter"
	"github.com/bhborges/family-tree-api/internal/app"
	"github.com/bhborges/family-tree-api/internal/port/rest"

	"go.uber.org/fx"
)

// APIModule wraps all logic related to the main family tree API.
func APIModule() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(adapter.NewPostgresRepository, fx.As(new(app.Repository))),
		),
		fx.Provide(
			fx.Annotate(app.NewApplication, fx.As(new(rest.Application))),
		),
		fx.Provide(rest.ProvideHTTPServer),
		fx.Invoke(rest.RegisterHandlers),
	)
}
