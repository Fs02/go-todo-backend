//go:build go1.18
// +build go1.18

package reltest

import (
	"github.com/go-rel/rel"
)

// EntityRepository mock
type EntityRepository[T any] struct {
	rel.EntityRepository[T]
	mock *Repository
}

var _ rel.EntityRepository[any] = (*EntityRepository[any])(nil)

// ExpectIterate apply mocks and expectations for Iterate
func (er *EntityRepository[T]) ExpectIterate(query rel.Query, options ...rel.IteratorOption) *EntityMockIterate[T] {
	return &EntityMockIterate[T]{er.mock.ExpectIterate(query, options...)}
}

// ExpectAggregate apply mocks and expectations for Aggregate
func (er *EntityRepository[T]) ExpectAggregate(query rel.Query, aggregate string, field string) *MockAggregate {
	return er.mock.ExpectAggregate(query, aggregate, field)
}

// ExpectCount apply mocks and expectations for Count
func (er *EntityRepository[T]) ExpectCount(collection string, queriers ...rel.Querier) *MockCount {
	return er.mock.ExpectCount(collection, queriers...)
}

// ExpectFind apply mocks and expectations for Find
func (er *EntityRepository[T]) ExpectFind(queriers ...rel.Querier) *EntityMockFind[T] {
	return &EntityMockFind[T]{er.mock.ExpectFind(queriers...)}
}

// ExpectFindAll apply mocks and expectations for FindAll
func (er *EntityRepository[T]) ExpectFindAll(queriers ...rel.Querier) *EntityMockFindAll[T] {
	return &EntityMockFindAll[T]{er.mock.ExpectFindAll(queriers...)}
}

// ExpectFindAndCountAll apply mocks and expectations for FindAndCountAll
func (er *EntityRepository[T]) ExpectFindAndCountAll(queriers ...rel.Querier) *EntityMockFindAndCountAll[T] {
	return &EntityMockFindAndCountAll[T]{er.mock.ExpectFindAndCountAll(queriers...)}
}

// ExpectInsert apply mocks and expectations for Insert
func (er *EntityRepository[T]) ExpectInsert(mutators ...rel.Mutator) *EntityMockMutate[T] {
	return &EntityMockMutate[T]{er.mock.ExpectInsert(mutators...)}
}

// ExpectInsertAll entities.
func (er *EntityRepository[T]) ExpectInsertAll() *EntityMockInsertAll[T] {
	return &EntityMockInsertAll[T]{er.mock.ExpectInsertAll()}
}

// ExpectUpdate apply mocks and expectations for Update
func (er *EntityRepository[T]) ExpectUpdate(mutators ...rel.Mutator) *EntityMockMutate[T] {
	return &EntityMockMutate[T]{er.mock.ExpectUpdate(mutators...)}
}

// ExpectDelete apply mocks and expectations for Delete
func (er *EntityRepository[T]) ExpectDelete(options ...rel.Mutator) *EntityMockDelete[T] {
	return &EntityMockDelete[T]{er.mock.ExpectDelete(options...)}
}

// ExpectDeleteAll apply mocks and expectations for DeleteAll
func (er *EntityRepository[T]) ExpectDeleteAll() *EntityMockDeleteAll[T] {
	return &EntityMockDeleteAll[T]{er.mock.ExpectDeleteAll()}
}

// ExpectPreload apply mocks and expectations for Preload
func (er *EntityRepository[T]) ExpectPreload(field string, queriers ...rel.Querier) *EntityMockPreload[T] {
	return &EntityMockPreload[T]{er.mock.ExpectPreload(field, queriers...)}
}

// ExpectPreload apply mocks and expectations for Preload
func (er *EntityRepository[T]) ExpectPreloadAll(field string, queriers ...rel.Querier) *EntityMockPreloadAll[T] {
	return &EntityMockPreloadAll[T]{er.mock.ExpectPreload(field, queriers...)}
}

// ExpectTransaction declare expectation inside transaction.
func (er *EntityRepository[T]) ExpectTransaction(fn func(*Repository)) {
	er.mock.ExpectTransaction(fn)
}

// AssertExpectations asserts that everything was in fact called as expected. Calls may have occurred in any order.
func (er *EntityRepository[T]) AssertExpectations(t TestingT) bool {
	t.Helper()
	return er.mock.AssertExpectations(t)
}

// WrapTransaction repository
func (er *EntityRepository[T]) WrapTransaction(repository *Repository) *EntityRepository[T] {
	return toEntityRepository[T](repository)
}

func NewEntityRepository[T any]() *EntityRepository[T] {
	repository := New()
	return &EntityRepository[T]{
		EntityRepository: rel.NewEntityRepository[T](repository),
		mock:             repository,
	}
}

func toEntityRepository[T any](repository *Repository) *EntityRepository[T] {
	return &EntityRepository[T]{
		EntityRepository: rel.NewEntityRepository[T](repository),
		mock:             repository,
	}
}
