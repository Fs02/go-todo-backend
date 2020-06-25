package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fs02/go-todo-backend/api/handler"
	"github.com/stretchr/testify/assert"
)

type pinger struct {
	err error
}

func (p pinger) Ping(ctx context.Context) error {
	return p.err
}

func TestHealthz_Show(t *testing.T) {
	tests := []struct {
		name     string
		pinger   handler.Pinger
		status   int
		path     string
		response string
	}{
		{
			name:     "all dependencies are healthy",
			pinger:   pinger{},
			status:   http.StatusOK,
			path:     "/",
			response: `[{"service": "test", "status": "UP"}]`,
		},
		{
			name:     "some dependencies are sick",
			pinger:   pinger{err: errors.New("service is down")},
			status:   http.StatusServiceUnavailable,
			path:     "/",
			response: `[{"service": "test", "status": "service is down"}]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _  = http.NewRequest("GET", test.path, nil)
				rr      = httptest.NewRecorder()
				handler = handler.NewHealthz()
			)

			handler.Add("test", test.pinger)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.JSONEq(t, test.response, rr.Body.String())
		})
	}
}
