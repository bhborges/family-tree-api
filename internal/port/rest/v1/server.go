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
	ListPeople(context.Context) ([]*domain.Person, error)
	GetPersonByID(context.Context, string) (*domain.Person, error)
	CreatePerson(context.Context, domain.Person) (string, error)
	CreatePeople(context.Context, []domain.Person) ([]string, error)
	UpdatePerson(context.Context, domain.Person) error
	DeletePerson(context.Context, string) error
	ListRelationships(context.Context) ([]*domain.Relationship, error)
	CreateRelationship(context.Context, domain.Relationship) (string, error)
	CreateRelationships(context.Context, []domain.Relationship) ([]string, error)
	UpdateRelationship(context.Context, domain.Relationship) error
	DeleteRelationship(context.Context, string) error
	BuildFamilyTree(context.Context, string) (*domain.FamilyTree, error)
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
		r.Use(http.FormatMiddleware)
		r.Use(http.SetContentTypeMiddleware)
		r.Route("/person", func(r chi.Router) {
			r.Get("/", http.WithAPM(h.apm, "/", h.ListPeople))
			r.Get("/{id}", http.WithAPM(h.apm, "/{id}", h.BuildFamilyTree))
			r.Post("/", http.WithAPM(h.apm, "/", h.CreatePerson))
			r.Patch("/", http.WithAPM(h.apm, "/", h.UpdatePerson))
			r.Delete("/{id}", http.WithAPM(h.apm, "/{id}", h.DeletePerson))
		})
		r.Route("/people", func(r chi.Router) {
			r.Post("/", http.WithAPM(h.apm, "/", h.CreatePeople))
		})
		r.Route("/relationship", func(r chi.Router) {
			r.Post("/", http.WithAPM(h.apm, "/", h.CreateRelationship))
			r.Put("/{id}", http.WithAPM(h.apm, "/{id}", h.UpdateRelationship))
			r.Delete("/{id}", http.WithAPM(h.apm, "/{id}", h.DeleteRelationship))
		})
		r.Route("/relationships", func(r chi.Router) {
			r.Get("/", http.WithAPM(h.apm, "/", h.ListRelationships))
			r.Post("/", http.WithAPM(h.apm, "/", h.CreateRelationships))
		})
	})
}
