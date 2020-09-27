package todos

import (
	"context"

	"github.com/Fs02/go-todo-backend/scores"
	"github.com/go-rel/rel"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "todos")))
)

//go:generate mockery --all --case=underscore --output todostest --outpkg todostest

// Service instance for todo's domain.
// Any operation done to any of object within this domain should use this service.
type Service interface {
	Search(ctx context.Context, todos *[]Todo, filter Filter) error
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, todo *Todo, changes rel.Changeset) error
	Delete(ctx context.Context, todo *Todo)
	Clear(ctx context.Context)
}

// beside embeding the struct, you can also declare the function directly on this struct.
// the advantage of embedding the struct is it allows spreading the implementation across multiple files.
type service struct {
	search
	create
	update
	delete
	clear
}

var _ Service = (*service)(nil)

// New Todos service.
func New(repository rel.Repository, scores scores.Service) Service {
	return service{
		search: search{repository: repository},
		create: create{repository: repository, scores: scores},
		update: update{repository: repository, scores: scores},
		delete: delete{repository: repository},
		clear:  clear{repository: repository},
	}
}
