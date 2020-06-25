package todos

import (
	"context"

	"github.com/Fs02/rel"
	"go.uber.org/zap"
)

type update struct {
	repository rel.Repository
}

func (u update) Update(ctx context.Context, todo *Todo, changes rel.Changeset) error {
	if err := todo.Validate(); err != nil {
		logger.Warn("validation error", zap.Error(err))
		return err
	}

	u.repository.MustUpdate(ctx, todo, changes)
	return nil
}
