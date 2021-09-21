package todos

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository, nil)
	)

	repository.ExpectDeleteAny(rel.From("todos")).Unsafe()

	assert.NotPanics(t, func() {
		service.Clear(ctx)
	})

	repository.AssertExpectations(t)
}
