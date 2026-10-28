package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fs02/go-todo-backend/api"
	"github.com/Fs02/go-todo-backend/scores"
	"github.com/Fs02/go-todo-backend/todos"
	"github.com/go-rel/rel"
	"github.com/go-rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestNewMux_Healthz(t *testing.T) {
	repository := reltest.New()

	mux := api.NewMux(repository)

	req, _ := http.NewRequest(http.MethodGet, "/healthz/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[{"service": "database", "status": "UP"}]`, rr.Body.String())
	repository.AssertExpectations(t)
}

func TestNewMux_Todos(t *testing.T) {
	repository := reltest.New()
	repository.ExpectFindAll(rel.Select().SortAsc("order")).Result([]todos.Todo{
		{ID: 1, Title: "Learn Go", Completed: false},
	})

	mux := api.NewMux(repository)

	req, _ := http.NewRequest(http.MethodGet, "/todos/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"title":"Learn Go"`)
	repository.AssertExpectations(t)
}

func TestNewMux_Score(t *testing.T) {
	repository := reltest.New()
	repository.ExpectFind().Result(scores.Score{ID: 1, TotalPoint: 42})

	mux := api.NewMux(repository)

	req, _ := http.NewRequest(http.MethodGet, "/score/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"total_point":42`)
	repository.AssertExpectations(t)
}

func TestNewMux_CORS(t *testing.T) {
	repository := reltest.New()
	mux := api.NewMux(repository)

	req, _ := http.NewRequest(http.MethodOptions, "/todos/", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}
