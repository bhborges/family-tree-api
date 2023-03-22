// Package family tree holds all modules related to it.
package familytree

import (
	"github.com/bhborges/family-tree-api/internal/port/rest"

	"go.uber.org/fx"
)

// APIModule wraps all logic related to the main family tree API.
func APIModule() fx.Option {
	return fx.Options(
		fx.Provide(rest.ProvideHTTPServer),
		fx.Invoke(rest.RegisterHandlers),
	)
}
