// Package rest implements an HTTP server.
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

// ListPeople returns a list of people.
func (h *HTTPServer) ListPeople(w http.ResponseWriter, r *http.Request) {
	p, err := h.application.ListPeople(r.Context())

	if errors.Is(err, app.ErrPersonNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%x", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error retrieving list of peoples from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, p)
}

// GetPersonByID returns a person.
func (h *HTTPServer) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.application.GetPersonByID(r.Context(), id)

	if errors.Is(err, app.ErrPersonNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%x", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error retrieving person from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, p)
}

// CreatePerson create a new person.
func (h *HTTPServer) CreatePerson(w http.ResponseWriter, r *http.Request) {
	p := domain.Person{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := h.application.CreatePerson(r.Context(), p)
	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error creating person from API", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, id)
}

func (h *HTTPServer) CreatePeople(w http.ResponseWriter, r *http.Request) {
	var people []domain.Person

	if err := json.NewDecoder(r.Body).Decode(&people); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var ids []string
	for _, p := range people {
		id, err := h.application.CreatePerson(r.Context(), p)
		if err != nil {
			newrelic.FromContext(r.Context()).NoticeError(err)
			h.log.Error("unexpected error creating person from API", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		ids = append(ids, id)
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ids)
}

// UpdatePerson update a person.
func (h *HTTPServer) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	p := domain.Person{}

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err := h.application.UpdatePerson(r.Context(), p)
	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error updating person from API", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
}

// DeletePerson delete a person.
func (h *HTTPServer) DeletePerson(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.application.DeletePerson(r.Context(), id)

	if errors.Is(err, app.ErrPersonNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%x", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error delete person from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
}

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

// CreateRelationShip create a new relationship.
func (h *HTTPServer) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	dr := domain.Relationship{}

	if err := json.NewDecoder(r.Body).Decode(&dr); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id, err := h.application.CreateRelationship(r.Context(), dr)
	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error creating relationship from API", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, id)
}

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
