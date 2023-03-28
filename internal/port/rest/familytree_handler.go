// Package rest implements an HTTP server.
package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bhborges/family-tree-api/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

// ListPeople returns a list of people.
func (h *HTTPServer) BuildFamilyTree(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	t, err := h.application.BuildFamilyTree(r.Context(), id)

	if errors.Is(err, app.ErrPersonNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%x", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error retrieving family tree from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, t)
}
