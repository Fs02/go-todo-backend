package handler

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		response string
	}{
		{
			name:     "message",
			data:     "lorem",
			response: `{"message":"lorem"}`,
		},
		{
			name:     "error",
			data:     errors.New("system error"),
			response: `{"error":"system error"}`,
		},
		{
			name:     "nil",
			data:     nil,
			response: ``,
		},
		{
			name: "struct",
			data: struct {
				ID int `json:"id"`
			}{ID: 1},
			response: `{"id":1}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				rr = httptest.NewRecorder()
			)

			render(rr, test.data, 200)
			if test.response != "" {
				assert.JSONEq(t, test.response, rr.Body.String())
			} else {
				assert.Equal(t, test.response, rr.Body.String())
			}
		})
	}
}
