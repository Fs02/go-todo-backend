package reltest

import (
	"context"
	"database/sql"
	"runtime"

	"github.com/go-rel/rel"
)

var (
	// ErrConnectionClosed is alias for sql.ErrConnDone.
	ErrConnectionClosed = sql.ErrConnDone
)

// Repository mock
type Repository struct {
	repo            rel.Repository
	iterate         iterate
	count           count
	aggregate       aggregate
	find            find
	findAll         findAll
	findAndCountAll findAndCountAll
	insert          mutate
	insertAll       insertAll
	update          mutate
	updateAny       updateAny
	delete          delete
	deleteAll       deleteAll
	deleteAny       deleteAny
	exec            exec
	preload         preload
	transaction     Assert
	ctxData         ctxData
}

var _ rel.Repository = (*Repository)(nil)

// Adapter provides a mock function with given fields:
func (r *Repository) Adapter(ctx context.Context) rel.Adapter {
	return r.repo.Adapter(ctx)
}

// Instrumentation provides a mock function with given fields: instrumenter
func (r *Repository) Instrumentation(instrumenter rel.Instrumenter) {
	r.repo.Instrumentation(instrumenter)
}

// Ping database.
func (r *Repository) Ping(ctx context.Context) error {
	return r.repo.Ping(ctx)
}

// Iterate through a collection of entities from database in batches.
// This function returns iterator that can be used to loop all entities.
// Limit, Offset and Sort query is automatically ignored.
func (r *Repository) Iterate(ctx context.Context, query rel.Query, options ...rel.IteratorOption) rel.Iterator {
	return r.iterate.execute(ctx, query, options...)
}

// ExpectIterate apply mocks and expectations for Iterate
func (r *Repository) ExpectIterate(query rel.Query, options ...rel.IteratorOption) *MockIterate {
	return r.iterate.register(r.ctxData, query, options...)
}

// Aggregate provides a mock function with given fields: query, aggregate, field
func (r *Repository) Aggregate(ctx context.Context, query rel.Query, aggregate string, field string) (int, error) {
	r.repo.Aggregate(ctx, query, aggregate, field)
	return r.aggregate.execute(ctx, query, aggregate, field)
}

// MustAggregate provides a mock function with given fields: query, aggregate, field
func (r *Repository) MustAggregate(ctx context.Context, query rel.Query, aggregate string, field string) int {
	result, err := r.Aggregate(ctx, query, aggregate, field)
	must(err)
	return result
}

// ExpectAggregate apply mocks and expectations for Aggregate
func (r *Repository) ExpectAggregate(query rel.Query, aggregate string, field string) *MockAggregate {
	return r.aggregate.register(r.ctxData, query, aggregate, field)
}

// Count provides a mock function with given fields: collection, queriers
func (r *Repository) Count(ctx context.Context, collection string, queriers ...rel.Querier) (int, error) {
	r.repo.Count(ctx, collection, queriers...)
	return r.count.execute(ctx, collection, queriers...)
}

// MustCount provides a mock function with given fields: collection, queriers
func (r *Repository) MustCount(ctx context.Context, collection string, queriers ...rel.Querier) int {
	count, err := r.Count(ctx, collection, queriers...)
	must(err)
	return count
}

// ExpectCount apply mocks and expectations for Count
func (r *Repository) ExpectCount(collection string, queriers ...rel.Querier) *MockCount {
	return r.count.register(r.ctxData, collection, queriers...)
}

// Find provides a mock function with given fields: entity, queriers
func (r *Repository) Find(ctx context.Context, entity any, queriers ...rel.Querier) error {
	r.repo.Find(ctx, entity, queriers...)
	return r.find.execute(ctx, entity, queriers...)
}

// MustFind provides a mock function with given fields: entity, queriers
func (r *Repository) MustFind(ctx context.Context, entity any, queriers ...rel.Querier) {
	must(r.Find(ctx, entity, queriers...))
}

// ExpectFind apply mocks and expectations for Find
func (r *Repository) ExpectFind(queriers ...rel.Querier) *MockFind {
	return r.find.register(r.ctxData, queriers...)
}

// FindAll provides a mock function with given fields: entities, queriers
func (r *Repository) FindAll(ctx context.Context, entities any, queriers ...rel.Querier) error {
	r.repo.FindAll(ctx, entities, queriers...)
	return r.findAll.execute(ctx, entities, queriers...)
}

// MustFindAll provides a mock function with given fields: entities, queriers
func (r *Repository) MustFindAll(ctx context.Context, entities any, queriers ...rel.Querier) {
	must(r.FindAll(ctx, entities, queriers...))
}

// ExpectFindAll apply mocks and expectations for FindAll
func (r *Repository) ExpectFindAll(queriers ...rel.Querier) *MockFindAll {
	return r.findAll.register(r.ctxData, queriers...)
}

// FindAndCountAll provides a mock function with given fields: entities, queriers
func (r *Repository) FindAndCountAll(ctx context.Context, entities any, queriers ...rel.Querier) (int, error) {
	r.repo.FindAndCountAll(ctx, entities, queriers...)
	return r.findAndCountAll.execute(ctx, entities, queriers...)
}

// MustFindAndCountAll provides a mock function with given fields: entities, queriers
func (r *Repository) MustFindAndCountAll(ctx context.Context, entities any, queriers ...rel.Querier) int {
	count, err := r.FindAndCountAll(ctx, entities, queriers...)
	must(err)
	return count
}

// ExpectFindAndCountAll apply mocks and expectations for FindAndCountAll
func (r *Repository) ExpectFindAndCountAll(queriers ...rel.Querier) *MockFindAndCountAll {
	return r.findAndCountAll.register(r.ctxData, queriers...)
}

// Insert provides a mock function with given fields: entity, mutators
func (r *Repository) Insert(ctx context.Context, entity any, mutators ...rel.Mutator) error {
	ret := r.insert.execute("Insert", ctx, entity, mutators...)
	r.repo.Insert(ctx, entity, mutators...)
	return ret
}

// MustInsert provides a mock function with given fields: entity, mutators
func (r *Repository) MustInsert(ctx context.Context, entity any, mutators ...rel.Mutator) {
	must(r.Insert(ctx, entity, mutators...))
}

// ExpectInsert apply mocks and expectations for Insert
func (r *Repository) ExpectInsert(mutators ...rel.Mutator) *MockMutate {
	return r.insert.register("Insert", r.ctxData, mutators...)
}

// InsertAll entities.
func (r *Repository) InsertAll(ctx context.Context, entities any, mutators ...rel.Mutator) error {
	r.repo.InsertAll(ctx, entities, mutators...)
	return r.insertAll.execute(ctx, entities)
}

// MustInsertAll entities.
func (r *Repository) MustInsertAll(ctx context.Context, entities any, mutators ...rel.Mutator) {
	must(r.InsertAll(ctx, entities, mutators...))
}

// ExpectInsertAll entities.
func (r *Repository) ExpectInsertAll() *MockInsertAll {
	return r.insertAll.register(r.ctxData)
}

// Update provides a mock function with given fields: entity, mutators
func (r *Repository) Update(ctx context.Context, entity any, mutators ...rel.Mutator) error {
	ret := r.update.execute("Update", ctx, entity, mutators...)
	if err := r.repo.Update(ctx, entity, mutators...); err != nil {
		return err
	}

	return ret
}

// MustUpdate provides a mock function with given fields: entity, mutators
func (r *Repository) MustUpdate(ctx context.Context, entity any, mutators ...rel.Mutator) {
	must(r.Update(ctx, entity, mutators...))
}

// ExpectUpdate apply mocks and expectations for Update
func (r *Repository) ExpectUpdate(mutators ...rel.Mutator) *MockMutate {
	return r.update.register("Update", r.ctxData, mutators...)
}

// UpdateAny provides a mock function with given fields: query
func (r *Repository) UpdateAny(ctx context.Context, query rel.Query, mutates ...rel.Mutate) (int, error) {
	return r.updateAny.execute(ctx, query, mutates...)
}

// MustUpdateAny provides a mock function with given fields: query
func (r *Repository) MustUpdateAny(ctx context.Context, query rel.Query, mutates ...rel.Mutate) int {
	r.repo.UpdateAny(ctx, query, mutates...)
	updatedCount, err := r.UpdateAny(ctx, query, mutates...)
	must(err)
	return updatedCount
}

// ExpectUpdateAny apply mocks and expectations for UpdateAny
func (r *Repository) ExpectUpdateAny(query rel.Query, mutates ...rel.Mutate) *MockUpdateAny {
	return r.updateAny.register(r.ctxData, query, mutates...)
}

// Delete provides a mock function with given fields: entity
func (r *Repository) Delete(ctx context.Context, entity any, options ...rel.Mutator) error {
	r.repo.Delete(ctx, entity)
	return r.delete.execute(ctx, entity, options...)
}

// MustDelete provides a mock function with given fields: entity
func (r *Repository) MustDelete(ctx context.Context, entity any, options ...rel.Mutator) {
	must(r.Delete(ctx, entity, options...))
}

// ExpectDelete apply mocks and expectations for Delete
func (r *Repository) ExpectDelete(options ...rel.Mutator) *MockDelete {
	return r.delete.register(r.ctxData, options...)
}

// DeleteAll provides DeleteAll mock function with given fields: entities
func (r *Repository) DeleteAll(ctx context.Context, entities any) error {
	r.repo.DeleteAll(ctx, entities)
	return r.deleteAll.execute(ctx, entities)
}

// MustDeleteAll provides a mock function with given fields: entity
func (r *Repository) MustDeleteAll(ctx context.Context, entity any) {
	must(r.DeleteAll(ctx, entity))
}

// ExpectDeleteAll apply mocks and expectations for DeleteAll
func (r *Repository) ExpectDeleteAll() *MockDeleteAll {
	return r.deleteAll.register(r.ctxData)
}

// DeleteAny provides a mock function with given fields: query
func (r *Repository) DeleteAny(ctx context.Context, query rel.Query) (int, error) {
	r.repo.DeleteAny(ctx, query)
	return r.deleteAny.execute(ctx, query)
}

// MustDeleteAny provides a mock function with given fields: query
func (r *Repository) MustDeleteAny(ctx context.Context, query rel.Query) int {
	deletedCount, err := r.DeleteAny(ctx, query)
	must(err)
	return deletedCount
}

// ExpectDeleteAny apply mocks and expectations for DeleteAny
func (r *Repository) ExpectDeleteAny(query rel.Query) *MockDeleteAny {
	return r.deleteAny.register(r.ctxData, query)
}

// Preload provides a mock function with given fields: entities, field, queriers
func (r *Repository) Preload(ctx context.Context, entities any, field string, queriers ...rel.Querier) error {
	return r.preload.execute(ctx, entities, field, queriers...)
}

// MustPreload provides a mock function with given fields: entities, field, queriers
func (r *Repository) MustPreload(ctx context.Context, entities any, field string, queriers ...rel.Querier) {
	must(r.Preload(ctx, entities, field, queriers...))
}

// ExpectPreload apply mocks and expectations for Preload
func (r *Repository) ExpectPreload(field string, queriers ...rel.Querier) *MockPreload {
	return r.preload.register(r.ctxData, field, queriers...)
}

// Exec raw statement.
// Returns last inserted id, rows affected and error.
func (r *Repository) Exec(ctx context.Context, statement string, args ...any) (int, int, error) {
	return r.exec.execute(ctx, statement, args)
}

// MustExec raw statement.
// Returns last inserted id, rows affected and error.
func (r *Repository) MustExec(ctx context.Context, statement string, args ...any) (int, int) {
	lastInsertedId, rowsAffected, err := r.Exec(ctx, statement, args...)
	must(err)
	return lastInsertedId, rowsAffected
}

// ExpectExec for mocking Exec
func (r *Repository) ExpectExec(statement string, args ...any) *MockExec {
	return r.exec.register(r.ctxData, statement, args...)
}

// Transaction provides a mock function with given fields: fn
func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	ctxData := fetchContext(ctx)
	r.transaction.call(ctx)

	var err error
	func() {
		defer func() {
			if p := recover(); p != nil {
				switch e := p.(type) {
				case runtime.Error:
					panic(e)
				case error:
					err = e
				default:
					panic(e)
				}
			}
		}()

		ctxData.txDepth++
		ctx = wrapContext(ctx, ctxData)
		err = fn(ctx)
	}()

	return err
}

// ExpectTransaction declare expectation inside transaction.
func (r *Repository) ExpectTransaction(fn func(*Repository)) {
	r.transaction.repeatability++

	r.ctxData.txDepth++
	fn(r)
	r.ctxData.txDepth--
}

// AssertExpectations asserts that everything was in fact called as expected. Calls may have occurred in any order.
func (r *Repository) AssertExpectations(t TestingT) bool {
	t.Helper()
	return r.iterate.assert(t) &&
		r.count.assert(t) &&
		r.aggregate.assert(t) &&
		r.find.assert(t) &&
		r.findAll.assert(t) &&
		r.findAndCountAll.assert(t) &&
		r.insert.assert(t) &&
		r.insertAll.assert(t) &&
		r.update.assert(t) &&
		r.updateAny.assert(t) &&
		r.delete.assert(t) &&
		r.deleteAll.assert(t) &&
		r.deleteAny.assert(t) &&
		r.exec.assert(t) &&
		r.preload.assert(t)
	// TODO: r.transaction.assert(t)
}

// New test repository.
func New() *Repository {
	return &Repository{
		repo: rel.New(&nopAdapter{}),
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
