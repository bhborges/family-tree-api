// Package rest implements an HTTP server.
package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

// GetPersonByID returns a person.
func (h *HTTPServer) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	p, err := h.application.GetPersonByID(r.Context(), id)
	if err != nil {
		h.log.Error("unexpected error retrieving person from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, p)
}
