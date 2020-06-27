package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Fs02/go-todo-backend/todos"
	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type ctx int

const (
	bodyKey ctx = 0
	loadKey ctx = 1
)

// Todos for todos endpoints.
type Todos struct {
	*chi.Mux
	repository rel.Repository
	todos      todos.Service
}

// Index handle GET /.
func (t Todos) Index(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		query  = r.URL.Query()
		result []todos.Todo
		filter = todos.Filter{
			Keyword: query.Get("keyword"),
		}
	)

	if str := query.Get("completed"); str != "" {
		completed := str == "true"
		filter.Completed = &completed
	}

	t.todos.Search(ctx, &result, filter)
	render(w, result, 200)
}

// Create handle POST /
func (t Todos) Create(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		todo todos.Todo
	)

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		logger.Warn("decode error", zap.Error(err))
		render(w, ErrBadRequest, 400)
		return
	}

	if err := t.todos.Create(ctx, &todo); err != nil {
		render(w, err, 422)
		return
	}

	w.Header().Set("Location", fmt.Sprint(r.RequestURI, "/", todo.ID))
	render(w, todo, 201)
}

// Show handle GET /{ID}
func (t Todos) Show(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		todo = ctx.Value(loadKey).(todos.Todo)
	)

	render(w, todo, 200)
}

// Update handle PATCH /{ID}
func (t Todos) Update(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		todo    = ctx.Value(loadKey).(todos.Todo)
		changes = rel.NewChangeset(&todo)
	)

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		logger.Warn("decode error", zap.Error(err))
		render(w, ErrBadRequest, 400)
		return
	}

	if err := t.todos.Update(ctx, &todo, changes); err != nil {
		render(w, err, 422)
		return
	}

	render(w, todo, 200)
}

// Destroy handle DELETE /{ID}
func (t Todos) Destroy(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		todo = ctx.Value(loadKey).(todos.Todo)
	)

	t.todos.Delete(ctx, &todo)
	render(w, nil, 204)
}

// Clear handle DELETE /
func (t Todos) Clear(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	t.todos.Clear(ctx)
	render(w, nil, 204)
}

// Load is middleware that loads todos to context.
func (t Todos) Load(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx   = r.Context()
			id, _ = strconv.Atoi(chi.URLParam(r, "ID"))
			todo  todos.Todo
		)

		if err := t.repository.Find(ctx, &todo, where.Eq("id", id)); err != nil {
			if errors.Is(err, rel.ErrNotFound) {
				render(w, err, 404)
				return
			}
			panic(err)
		}

		ctx = context.WithValue(ctx, loadKey, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewTodos handler.
func NewTodos(repository rel.Repository, todos todos.Service) Todos {
	h := Todos{
		Mux:        chi.NewMux(),
		repository: repository,
		todos:      todos,
	}

	h.Get("/", h.Index)
	h.Post("/", h.Create)
	h.With(h.Load).Get("/{ID}", h.Show)
	h.With(h.Load).Patch("/{ID}", h.Update)
	h.With(h.Load).Delete("/{ID}", h.Destroy)
	h.Delete("/", h.Clear)

	return h
}
