package todos

import (
	"context"

	"github.com/Fs02/rel"
	"go.uber.org/zap"
)

type create struct {
	repository rel.Repository
}

func (c create) Create(ctx context.Context, todo *Todo) error {
	if err := todo.Validate(); err != nil {
		logger.Warn("validation error", zap.Error(err))
		return err
	}

	c.repository.MustInsert(ctx, todo)
	return nil
}
