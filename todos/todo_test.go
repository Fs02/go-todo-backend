package todos

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	TodoURLPrefix = "http://localhost:3000/"
}

func TestTodo_Validate(t *testing.T) {
	var todo Todo

	t.Run("title is blank", func(t *testing.T) {
		assert.Equal(t, ErrTodoTitleBlank, todo.Validate())
	})

	t.Run("valid", func(t *testing.T) {
		todo.Title = "Sleep"
		assert.Nil(t, todo.Validate())
	})
}

func TestTodo_MarshalJSON(t *testing.T) {
	var (
		todo = Todo{
			ID:        1,
			Title:     "Sleep",
			Completed: true,
		}
		encoded, err = json.Marshal(todo)
	)

	assert.Nil(t, err)
	assert.JSONEq(t, `{
		"id": 1,
		"title": "Sleep",
		"completed": true,
		"order": 0,
		"url": "http://localhost:3000/1",
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z"
	}`, string(encoded))
}
