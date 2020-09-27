package scores

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestEarn(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		name       = "todo completed"
		count      = 1
	)

	repository.ExpectTransaction(func(repository *reltest.Repository) {
		repository.ExpectFind(rel.ForUpdate()).Result(Score{ID: 1, TotalPoint: 10})
		repository.ExpectUpdate().For(&Score{ID: 1, TotalPoint: 11})
		repository.ExpectInsert().For(&Point{Name: name, Count: count, ScoreID: 1})
	})

	assert.Nil(t, service.Earn(ctx, name, count))
	repository.AssertExpectations(t)
}

func TestEarn_insertScore(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		name       = "todo completed"
		count      = 1
	)

	repository.ExpectTransaction(func(repository *reltest.Repository) {
		repository.ExpectFind(rel.ForUpdate()).NotFound()
		repository.ExpectInsert().For(&Score{TotalPoint: 1})
		repository.ExpectInsert().For(&Point{Name: name, Count: count, ScoreID: 1})
	})

	assert.Nil(t, service.Earn(ctx, name, count))
	repository.AssertExpectations(t)
}

func TestEarn_findError(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository)
		name       = "todo completed"
		count      = 1
	)

	repository.ExpectTransaction(func(repository *reltest.Repository) {
		repository.ExpectFind(rel.ForUpdate()).ConnectionClosed()
	})

	assert.Equal(t, reltest.ErrConnectionClosed, service.Earn(ctx, name, count))

	repository.AssertExpectations(t)
}
