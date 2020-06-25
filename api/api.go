package api

import (
	"github.com/Fs02/rel"
	"github.com/Fs02/go-todo-backend/api/handler"
	"github.com/Fs02/go-todo-backend/todos"
	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"
	"github.com/goware/cors"
)

// NewMux api.
func NewMux(repository rel.Repository) *chi.Mux {
	var (
		mux            = chi.NewMux()
		todos          = todos.New(repository)
		healthzHandler = handler.NewHealthz()
		todosHandler   = handler.NewTodos(repository, todos)
	)

	healthzHandler.Add("database", repository)

	mux.Use(chimid.RequestID)
	mux.Use(chimid.RealIP)
	mux.Use(chimid.Recoverer)
	mux.Use(cors.AllowAll().Handler)

	mux.Mount("/healthz", healthzHandler)
	mux.Mount("/todos", todosHandler)

	return mux
}
