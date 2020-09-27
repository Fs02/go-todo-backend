package scores

import (
	"context"
	"errors"

	"github.com/go-rel/rel"
)

type earn struct {
	repository rel.Repository
}

func (e earn) Earn(ctx context.Context, name string, count int) error {
	var (
		score Score
	)

	return e.repository.Transaction(ctx, func(ctx context.Context) error {
		// for simplicity, assumes only one user, so there's only one score and always retrieve the first one.
		// this will probably lock the entire table since there's no where clause provided, but it's find since we assume only one user.
		if err := e.repository.Find(ctx, &score, rel.ForUpdate()); err != nil {
			if !errors.Is(err, rel.ErrNotFound) {
				// unexpected error.
				return err
			}

			score.TotalPoint = count
			e.repository.MustInsert(ctx, &score)
		} else {
			score.TotalPoint += count
			e.repository.Update(ctx, &score)
		}

		// insert point history.
		e.repository.MustInsert(ctx, &Point{Name: name, Count: count, ScoreID: score.ID})
		return nil
	})
}
