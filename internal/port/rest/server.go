// Package rest family tree router
package rest

import (
	"context"

	"github.com/bhborges/family-tree-api/internal/domain"
	"github.com/bhborges/family-tree-api/pkg/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

// HTTPServer definition.
type HTTPServer struct {
	router *chi.Mux
	log    *zap.Logger
	apm    *newrelic.Application

	application Application
}

// Application specifies the signature of Application.
type Application interface {
	GetPersonByID(context.Context, string) (*domain.Person, error)
}

// ProvideHTTPServer returns a new instance of an HTTP server.
func ProvideHTTPServer(
	r *chi.Mux, l *zap.Logger,
	apm *newrelic.Application,
	application Application,
) *HTTPServer {
	return &HTTPServer{
		router:      r,
		log:         l,
		apm:         apm,
		application: application,
	}
}

// RegisterHandlers registers all handlers.
func RegisterHandlers(h *HTTPServer) {
	h.router.Route("/familytree", func(r chi.Router) {
		r.Route("/person", func(r chi.Router) {
			r.Get("/{id}", http.WithAPM(h.apm, "/", h.GetPersonByID))
		})

	})
}
