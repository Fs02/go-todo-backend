package todos

import (
	"context"

	"github.com/go-rel/rel"
)

// Filter for search.
type Filter struct {
	Keyword   string
	Completed *bool
}

type search struct {
	repository rel.Repository
}

func (s search) Search(ctx context.Context, todos *[]Todo, filter Filter) error {
	var (
		query = rel.Select().SortAsc("order")
	)

	if filter.Keyword != "" {
		query = query.Where(rel.Like("title", "%"+filter.Keyword+"%"))
	}

	if filter.Completed != nil {
		query = query.Where(rel.Eq("completed", *filter.Completed))
	}

	s.repository.MustFindAll(ctx, todos, query)
	return nil
}
