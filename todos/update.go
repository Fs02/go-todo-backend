package todos

import (
	"context"

	"github.com/Fs02/go-todo-backend/scores"
	"github.com/go-rel/rel"
	"go.uber.org/zap"
)

type update struct {
	repository rel.Repository
	scores     scores.Service
}

func (u update) Update(ctx context.Context, todo *Todo, changes rel.Changeset) error {
	if err := todo.Validate(); err != nil {
		logger.Warn("validation error", zap.Error(err))
		return err
	}

	// update score if completed is changed.
	if changes.FieldChanged("completed") {
		return u.repository.Transaction(ctx, func(ctx context.Context) error {
			u.repository.MustUpdate(ctx, todo, changes)

			if todo.Completed {
				return u.scores.Earn(ctx, "todo completed", 1)
			}

			return u.scores.Earn(ctx, "todo uncompleted", -2)
		})
	}

	u.repository.MustUpdate(ctx, todo, changes)
	return nil
}
