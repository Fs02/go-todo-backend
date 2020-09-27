package handler

import (
	"net/http"

	"github.com/Fs02/go-todo-backend/scores"
	"github.com/go-chi/chi"
	"github.com/go-rel/rel"
)

// Score for score endpoints.
type Score struct {
	*chi.Mux
	repository rel.Repository
}

// Index handle GET /
func (s Score) Index(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		result scores.Score
	)

	s.repository.Find(ctx, &result)
	render(w, result, 200)
}

// Points handle Get /points
func (s Score) Points(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		result []scores.Point
	)

	s.repository.FindAll(ctx, &result)
	render(w, result, 200)
}

// NewScore handler.
func NewScore(repository rel.Repository) Score {
	h := Score{
		Mux:        chi.NewMux(),
		repository: repository,
	}

	h.Get("/", h.Index)
	h.Get("/points", h.Points)

	return h
}
