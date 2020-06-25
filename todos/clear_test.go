package todos

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
	)

	repository.ExpectDeleteAll(rel.From("todos")).Unsafe()

	assert.NotPanics(t, func() {
		service.Clear(ctx)
	})

	repository.AssertExpectations(t)
}
