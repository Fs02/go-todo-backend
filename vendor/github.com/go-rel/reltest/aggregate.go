package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type aggregate []*MockAggregate

func (a *aggregate) register(ctxData ctxData, query rel.Query, aggregate string, field string) *MockAggregate {
	ma := &MockAggregate{
		assert:       &Assert{ctxData: ctxData},
		argQuery:     query,
		argAggregate: aggregate,
		argField:     field,
	}
	*a = append(*a, ma)
	return ma
}

func (a aggregate) execute(ctx context.Context, query rel.Query, aggregate string, field string) (int, error) {
	for _, ma := range a {
		if matchQuery(ma.argQuery, query) &&
			ma.argAggregate == aggregate &&
			ma.argField == field &&
			ma.assert.call(ctx) {
			return ma.retCount, ma.retError
		}
	}

	ma := &MockAggregate{
		assert:       &Assert{ctxData: fetchContext(ctx)},
		argQuery:     query,
		argAggregate: aggregate,
		argField:     field,
	}
	panic(failExecuteMessage(ma, a))
}

func (a *aggregate) assert(t T) bool {
	t.Helper()
	for _, ma := range *a {
		if !ma.assert.assert(t, ma) {
			return false
		}
	}

	*a = nil
	return true
}

// MockAggregate asserts and simulate UpdateAny function for test.
type MockAggregate struct {
	assert       *Assert
	argQuery     rel.Query
	argAggregate string
	argField     string
	retCount     int
	retError     error
}

// Result sets the result of this query.
func (ma *MockAggregate) Result(count int) *Assert {
	ma.retCount = count
	return ma.assert
}

// Error sets error to be returned.
func (ma *MockAggregate) Error(err error) *Assert {
	ma.retError = err
	return ma.assert
}

// ConnectionClosed sets this error to be returned.
func (ma *MockAggregate) ConnectionClosed() *Assert {
	ma.Error(ErrConnectionClosed)
	return ma.assert
}

// String representation of mocked call.
func (ma MockAggregate) String() string {
	return ma.assert.sprintf(`Aggregate(ctx, %s, "%s", "%s")`, ma.argQuery, ma.argAggregate, ma.argField)
}

// ExpectString representation of mocked call.
func (ma MockAggregate) ExpectString() string {
	return ma.assert.sprintf(`ExpectAggregate(%s, "%s", "%s")`, ma.argQuery, ma.argAggregate, ma.argField)
}
