package todos

import (
	"context"

	"github.com/Fs02/rel"
)

type delete struct {
	repository rel.Repository
}

func (d delete) Delete(ctx context.Context, todo *Todo) {
	d.repository.MustDelete(ctx, todo)
}
