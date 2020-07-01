package todos

import (
	"context"
	"testing"

	"github.com/Fs02/go-todo-backend/scores/scorestest"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{Title: "Sleep"}
	)

	repository.ExpectInsert().For(&todo)

	assert.Nil(t, service.Create(ctx, &todo))
	assert.NotEmpty(t, todo.ID)

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}

func TestCreate_completed(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{Title: "Sleep", Completed: true}
	)

	repository.ExpectTransaction(func(repository *reltest.Repository) {
		scores.On("Earn", mock.Anything, "todo completed", 1).Return(nil)
		repository.ExpectInsert().For(&todo)
	})

	assert.Nil(t, service.Create(ctx, &todo))
	assert.NotEmpty(t, todo.ID)

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}

func TestCreate_validateError(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{Title: ""}
	)

	assert.Equal(t, ErrTodoTitleBlank, service.Create(ctx, &todo))

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}
