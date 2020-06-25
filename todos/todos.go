package todos

import (
	"context"

	"github.com/Fs02/rel"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "todos")))
)

//go:generate mockery -all -case=underscore -output todostest -outpkg todostest

// Todos service.
type Todos interface {
	Search(ctx context.Context, todos *[]Todo, filter Filter) error
	Create(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, todo *Todo, changes rel.Changeset) error
	Delete(ctx context.Context, todo *Todo)
	Clear(ctx context.Context)
}

// beside embeding the struct, you can also declare the function directly on this struct.
// the advantage of embedding the struct is it allows spreading the implementation across multiple files.
type todos struct {
	search
	create
	update
	delete
	clear
}

var _ Todos = (*todos)(nil)

// New Todos service.
func New(repository rel.Repository) Todos {
	return todos{
		search: search{repository: repository},
		create: create{repository: repository},
		update: update{repository: repository},
		delete: delete{repository: repository},
		clear:  clear{repository: repository},
	}
}
