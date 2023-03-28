package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bhborges/family-tree-api/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

// BuildFamilyTree returns a family tree.
func (h *HTTPServer) BuildFamilyTree(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	t, err := h.application.BuildFamilyTree(r.Context(), id)

	if errors.Is(err, app.ErrPersonNotFound) {
		render.Status(r, http.StatusNotFound)
		render.PlainText(w, r, fmt.Sprintf("%s", app.ErrPersonNotFound))
	}

	if err != nil {
		newrelic.FromContext(r.Context()).NoticeError(err)
		h.log.Error("unexpected error retrieving family tree from API server", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	render.Status(r, http.StatusOK)

	switch r.Header.Get("Accept") {
	case "application/xml":
		render.XML(w, r, t)
	case "application/octet-stream":
		bytes, _ := json.Marshal(t)
		render.Data(w, r, bytes)
	default:
		render.JSON(w, r, t)
	}
}
