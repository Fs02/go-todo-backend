package todos

import (
	"context"
	"testing"

	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		todo       = Todo{ID: 1, Title: "Sleep"}
	)

	repository.ExpectDelete().ForType("todos.Todo")

	assert.NotPanics(t, func() {
		service.Delete(ctx, &todo)
	})

	repository.AssertExpectations(t)
}
