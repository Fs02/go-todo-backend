package todos

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	var (
		ctx        = context.TODO()
		repository = reltest.New()
		service    = New(repository, nil)
		todos      []Todo
		completed  = false
		filter     = Filter{Keyword: "Sleep", Completed: &completed}
		result     = []Todo{{ID: 1, Title: "Sleep"}}
	)

	repository.ExpectFindAll(
		rel.Select().SortAsc("order").Where(rel.Like("title", "%Sleep%").AndEq("completed", false)),
	).Result(result)

	assert.NotPanics(t, func() {
		service.Search(ctx, &todos, filter)
		assert.Equal(t, result, todos)
	})

	repository.AssertExpectations(t)
}
