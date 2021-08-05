package scores

import (
	"context"

	"github.com/go-rel/rel"
)

//go:generate mockery --name=Service --case=underscore --output scorestest --outpkg scorestest

// Service instance for todo's domain.
// Any operation done to any of object within this domain should use this service.
type Service interface {
	Earn(ctx context.Context, name string, count int) error
}

// beside embeding the struct, you can also declare the function directly on this struct.
// the advantage of embedding the struct is it allows spreading the implementation across multiple files.
type service struct {
	earn
}

var _ Service = (*service)(nil)

// New Scores service.
func New(repository rel.Repository) Service {
	return service{
		earn: earn{repository: repository},
	}
}
