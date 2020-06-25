package todos

import (
	"context"
	"testing"

	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		todo       = Todo{Title: "Sleep"}
	)

	repository.ExpectInsert().For(&todo)

	assert.Nil(t, service.Create(ctx, &todo))
	assert.NotEmpty(t, todo.ID)
	repository.AssertExpectations(t)
}

func TestCreate_ValidateError(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		todo       = Todo{Title: ""}
	)

	assert.Equal(t, ErrTodoTitleBlank, service.Create(ctx, &todo))
	repository.AssertExpectations(t)
}
