package todos

import (
	"context"
	"testing"

	"github.com/Fs02/go-todo-backend/scores/scorestest"
	"github.com/Fs02/rel"

	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdate(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{ID: 1, Title: "Sleep"}
		changes    = rel.NewChangeset(&todo)
	)

	todo.Title = "Wake up"

	repository.ExpectUpdate(changes).ForType("todos.Todo")

	assert.Nil(t, service.Update(ctx, &todo, changes))
	assert.NotEmpty(t, todo.ID)

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}

func TestUpdate_completed(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{ID: 1, Title: "Sleep"}
		changes    = rel.NewChangeset(&todo)
	)

	todo.Completed = true

	repository.ExpectTransaction(func(repository *reltest.Repository) {
		scores.On("Earn", mock.Anything, "todo completed", 1).Return(nil)
		repository.ExpectUpdate(changes).ForType("todos.Todo")
	})

	assert.Nil(t, service.Update(ctx, &todo, changes))
	assert.NotEmpty(t, todo.ID)

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}

func TestUpdate_uncompleted(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{ID: 1, Title: "Sleep", Completed: true}
		changes    = rel.NewChangeset(&todo)
	)

	todo.Completed = false

	repository.ExpectTransaction(func(repository *reltest.Repository) {
		scores.On("Earn", mock.Anything, "todo uncompleted", -2).Return(nil)
		repository.ExpectUpdate(changes).ForType("todos.Todo")
	})

	assert.Nil(t, service.Update(ctx, &todo, changes))
	assert.NotEmpty(t, todo.ID)

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}

func TestUpdate_validateError(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		scores     = &scorestest.Service{}
		service    = New(repository, scores)
		todo       = Todo{ID: 1, Title: "Sleep"}
		changes    = rel.NewChangeset(&todo)
	)

	todo.Title = ""

	assert.Equal(t, ErrTodoTitleBlank, service.Update(ctx, &todo, changes))

	repository.AssertExpectations(t)
	scores.AssertExpectations(t)
}
