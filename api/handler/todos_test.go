package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Fs02/go-todo-backend/api/handler"
	"github.com/Fs02/go-todo-backend/todos"
	"github.com/Fs02/go-todo-backend/todos/todostest"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestTodos_Show(t *testing.T) {
	var (
		trueb = true
	)

	tests := []struct {
		name            string
		status          int
		path            string
		response        string
		mockTodosSearch func(todos *todostest.Todos)
	}{
		{
			name:     "ok",
			status:   http.StatusOK,
			path:     "/",
			response: `[{"id":1, "title":"Sleep", "completed":false, "order":0, "url":"todos/1"}]`,
			mockTodosSearch: todostest.MockTodosSearch(
				[]todos.Todo{{ID: 1, Title: "Sleep"}},
				todos.Filter{},
				nil,
			),
		},
		{
			name:     "with keyword and filter completed",
			status:   http.StatusOK,
			path:     "/?keyword=Wake&completed=true",
			response: `[{"id":2, "title":"Wake", "completed":true, "order":0, "url":"todos/2"}]`,
			mockTodosSearch: todostest.MockTodosSearch(
				[]todos.Todo{{ID: 2, Title: "Wake", Completed: true}},
				todos.Filter{Keyword: "Wake", Completed: &trueb},
				nil,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _     = http.NewRequest("GET", test.path, nil)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Todos{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.MockTodos(todos, test.mockTodosSearch)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.JSONEq(t, test.response, rr.Body.String())

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Create(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		path            string
		payload         string
		response        string
		location        string
		mockTodosCreate func(todos *todostest.Todos)
	}{
		{
			name:     "created",
			status:   http.StatusCreated,
			path:     "/",
			payload:  `{"title": "Sleep"}`,
			response: `{"id":1, "title":"Sleep", "completed":false, "order":0, "url":"todos/1"}`,
			location: "/1",
			mockTodosCreate: todostest.MockTodosCreate(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
		},
		{
			name:     "validation error",
			status:   http.StatusUnprocessableEntity,
			path:     "/",
			payload:  `{"title": ""}`,
			response: `{"error":"Title can't be blank"}`,
			mockTodosCreate: todostest.MockTodosCreate(
				todos.Todo{Title: "Sleep"},
				todos.ErrTodoTitleBlank,
			),
		},
		{
			name:     "bad request",
			status:   http.StatusBadRequest,
			path:     "/",
			payload:  ``,
			response: `{"error":"Bad Request"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				body       = strings.NewReader(test.payload)
				req, _     = http.NewRequest("POST", test.path, body)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Todos{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.MockTodos(todos, test.mockTodosCreate)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.Equal(t, test.location, rr.Header().Get("Location"))
			assert.JSONEq(t, test.response, rr.Body.String())

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Clear(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		path           string
		response       string
		mockTodosClear func(todos *todostest.Todos)
	}{
		{
			name:           "created",
			status:         http.StatusNoContent,
			path:           "/",
			response:       "",
			mockTodosClear: todostest.MockTodosClear(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _     = http.NewRequest("DELETE", test.path, nil)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Todos{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.MockTodos(todos, test.mockTodosClear)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			if test.response != "" {
				assert.JSONEq(t, test.response, rr.Body.String())
			} else {
				assert.Equal(t, "", rr.Body.String())
			}

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}
