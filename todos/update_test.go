package todos

import (
	"context"
	"testing"

	"github.com/Fs02/rel"

	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		todo       = Todo{ID: 1, Title: "Sleep"}
		changes    = rel.NewChangeset(&todo)
	)

	todo.Title = "Wake up"

	repository.ExpectUpdate(changes).ForType("todos.Todo")

	assert.Nil(t, service.Update(ctx, &todo, changes))
	assert.NotEmpty(t, todo.ID)
	repository.AssertExpectations(t)
}

func TestUpdate_ValidateError(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		todo       = Todo{ID: 1, Title: "Sleep"}
		changes    = rel.NewChangeset(&todo)
	)

	todo.Title = ""

	assert.Equal(t, ErrTodoTitleBlank, service.Update(ctx, &todo, changes))
	repository.AssertExpectations(t)
}
