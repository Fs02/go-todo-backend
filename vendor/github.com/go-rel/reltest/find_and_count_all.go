package reltest

import (
	"context"
	"reflect"

	"github.com/go-rel/rel"
)

type findAndCountAll []*MockFindAndCountAll

func (fca *findAndCountAll) register(ctxData ctxData, queriers ...rel.Querier) *MockFindAndCountAll {
	mfca := &MockFindAndCountAll{
		assert:   &Assert{ctxData: ctxData, repeatability: 1},
		argQuery: rel.Build("", queriers...),
	}
	*fca = append(*fca, mfca)
	return mfca
}

func (fca findAndCountAll) execute(ctx context.Context, entities any, queriers ...rel.Querier) (int, error) {
	query := rel.Build("", queriers...)
	for _, mfca := range fca {
		if matchQuery(mfca.argQuery, query) &&
			mfca.assert.call(ctx) {
			if mfca.argEntities != nil {
				reflect.ValueOf(entities).Elem().Set(reflect.ValueOf(mfca.argEntities))
			}

			return mfca.retCount, mfca.retError
		}
	}

	mfca := &MockFindAndCountAll{
		assert:      &Assert{ctxData: fetchContext(ctx)},
		argQuery:    query,
		argEntities: entities,
	}
	panic(failExecuteMessage(mfca, fca))
}

func (fca *findAndCountAll) assert(t TestingT) bool {
	t.Helper()
	for _, mfca := range *fca {
		if !mfca.assert.assert(t, mfca) {
			return false
		}
	}

	*fca = nil
	return true
}

// MockFindAndCountAll asserts and simulate find and count all function for test.
type MockFindAndCountAll struct {
	assert      *Assert
	argQuery    rel.Query
	argEntities any
	retCount    int
	retError    error
}

// Result sets the result of this query.
func (mfca *MockFindAndCountAll) Result(result any, count int) *Assert {
	if mfca.argQuery.Table == "" {
		mfca.argQuery.Table = rel.NewCollection(result, true).Table()
	}
	mfca.argEntities = result
	mfca.retCount = count
	return mfca.assert
}

// Error sets error to be returned.
func (mfca *MockFindAndCountAll) Error(err error) *Assert {
	mfca.retError = err
	return mfca.assert
}

// ConnectionClosed sets this error to be returned.
func (mfca *MockFindAndCountAll) ConnectionClosed() *Assert {
	return mfca.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (mfca MockFindAndCountAll) String() string {
	return mfca.assert.sprintf("FindAndCountAll(ctx, <Any>, %s)", mfca.argQuery)
}

// ExpectString representation of mocked call.
func (mfca MockFindAndCountAll) ExpectString() string {
	return mfca.assert.sprintf("ExpectFindAndCountAll(%s)", mfca.argQuery)
}
