package reltest

import (
	"context"
	"reflect"

	"github.com/go-rel/rel"
)

type findAll []*MockFindAll

func (fa *findAll) register(ctxData ctxData, queriers ...rel.Querier) *MockFindAll {
	mfa := &MockFindAll{
		assert:   &Assert{ctxData: ctxData, repeatability: 1},
		argQuery: rel.Build("", queriers...),
	}
	*fa = append(*fa, mfa)
	return mfa
}

func (fa findAll) execute(ctx context.Context, entities any, queriers ...rel.Querier) error {
	query := rel.Build("", queriers...)
	for _, mfa := range fa {
		if matchQuery(mfa.argQuery, query) &&
			mfa.assert.call(ctx) {
			if mfa.argEntities != nil {
				reflect.ValueOf(entities).Elem().Set(reflect.ValueOf(mfa.argEntities))
			}

			return mfa.retError
		}
	}

	mfa := &MockFindAll{
		assert:      &Assert{ctxData: fetchContext(ctx)},
		argQuery:    query,
		argEntities: entities,
	}
	panic(failExecuteMessage(mfa, fa))
}

func (fa *findAll) assert(t TestingT) bool {
	t.Helper()
	for _, mfa := range *fa {
		if !mfa.assert.assert(t, mfa) {
			return false
		}
	}

	*fa = nil
	return true
}

// MockFindAll asserts and simulate find all function for test.
type MockFindAll struct {
	assert      *Assert
	argQuery    rel.Query
	argEntities any
	retError    error
}

// Result sets the result of this query.
func (mfa *MockFindAll) Result(result any) *Assert {
	if mfa.argQuery.Table == "" {
		mfa.argQuery.Table = rel.NewCollection(result, true).Table()
	}
	mfa.argEntities = result
	return mfa.assert
}

// Error sets error to be returned.
func (mfa *MockFindAll) Error(err error) *Assert {
	mfa.retError = err
	return mfa.assert
}

// ConnectionClosed sets this error to be returned.
func (mfa *MockFindAll) ConnectionClosed() *Assert {
	return mfa.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (mfa MockFindAll) String() string {
	return mfa.assert.sprintf("FindAll(ctx, <Any>, %s)", mfa.argQuery)
}

// ExpectString representation of mocked call.
func (mfa MockFindAll) ExpectString() string {
	return mfa.assert.sprintf("ExpectFindAll(%s)", mfa.argQuery)
}
