package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bhborges/family-tree-api/internal/app"
	"github.com/bhborges/family-tree-api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

// ListRelationships list all relationships
func (h *HTTPServer) ListRelationships(w http.ResponseWriter, r *http.Request) {
	p, err := h.application.ListRelationships(r.Context())

	if errors.Is(err, app.ErrRelationshipNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrRelationshipNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error retrieving list of relatioships from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)

	switch r.Header.Get("Accept") {
	case "application/xml":
		render.XML(w, r, p)
	case "application/octet-stream":
		bytes, _ := json.Marshal(p)
		render.Data(w, r, bytes)
	default:
		render.JSON(w, r, p)
	}

}

// UpdateRelationship updates an existing relationship.
func (h *HTTPServer) UpdateRelationship(w http.ResponseWriter, r *http.Request) {
	dr := domain.Relationship{}

	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := h.application.UpdateRelationship(r.Context(), dr); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error updating relationship from API", zap.Error(err))

		if errors.Is(err, app.ErrRelationshipNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteRelationship deletes a relationship.
func (h *HTTPServer) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.application.DeleteRelationship(r.Context(), id); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error deleting relationship from API", zap.Error(err))

		if errors.Is(err, app.ErrRelationshipNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateRelationship create a new relationship.
func (h *HTTPServer) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	dr := domain.Relationship{}

	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := h.application.CreateRelationship(r.Context(), dr)

	if errors.Is(err, app.ErrIncestuousOffspring) {
		render.Status(r, http.StatusUnprocessableEntity)
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrIncestuousOffspring))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error creating relationship from API", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, id)
}

// CreateRelationships creates a new batch of relationships.
func (h *HTTPServer) CreateRelationships(w http.ResponseWriter, r *http.Request) {
	drs := []domain.Relationship{}

	if err := json.NewDecoder(r.Body).Decode(&drs); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ids, err := h.application.CreateRelationships(r.Context(), drs)
	if errors.Is(err, app.ErrIncestuousOffspring) {
		render.Status(r, http.StatusUnprocessableEntity)
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrIncestuousOffspring))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error creating relationships from API", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ids)
}
