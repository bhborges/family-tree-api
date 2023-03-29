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
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error retrieving list of peoples from API server", zap.Error(err))
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

// GetPersonByID returns a person.
func (h *HTTPServer) GetPersonByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.application.GetPersonByID(r.Context(), id)

	if errors.Is(err, app.ErrPersonNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrPersonNotFound))
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

// CreatePeople creates a new batch of people.
func (h *HTTPServer) CreatePeople(w http.ResponseWriter, r *http.Request) {
	var people []domain.Person

	if err := json.NewDecoder(r.Body).Decode(&people); err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error decoding data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ids := make([]string, 0, len(people))

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
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error delete person from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)
}
