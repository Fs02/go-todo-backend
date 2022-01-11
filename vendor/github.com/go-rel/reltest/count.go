package reltest

import (
	"context"

	"github.com/go-rel/rel"
)

type count []*MockCount

func (c *count) register(ctxData ctxData, collection string, queriers ...rel.Querier) *MockCount {
	mc := &MockCount{
		assert:        &Assert{ctxData: ctxData, repeatability: 1},
		argCollection: collection,
		argQuery:      rel.Build(collection, queriers...),
	}
	*c = append(*c, mc)
	return mc
}

func (c count) execute(ctx context.Context, collection string, queriers ...rel.Querier) (int, error) {
	query := rel.Build(collection, queriers...)
	for _, mc := range c {
		if mc.argCollection == collection &&
			matchQuery(mc.argQuery, query) &&
			mc.assert.call(ctx) {
			return mc.retCount, mc.retError
		}
	}

	mc := &MockCount{
		assert:        &Assert{ctxData: fetchContext(ctx)},
		argCollection: collection,
		argQuery:      query,
	}
	panic(failExecuteMessage(mc, c))
}

func (c *count) assert(t T) bool {
	t.Helper()
	for _, mc := range *c {
		if !mc.assert.assert(t, mc) {
			return false
		}
	}

	*c = nil
	return true
}

// MockCount asserts and simulate UpdateAny function for test.
type MockCount struct {
	assert        *Assert
	argCollection string
	argQuery      rel.Query
	retCount      int
	retError      error
}

// Result sets the result of this query.
func (mc *MockCount) Result(count int) *Assert {
	mc.retCount = count
	return mc.assert
}

// Error sets error to be returned.
func (mc *MockCount) Error(err error) *Assert {
	mc.retError = err
	return mc.assert
}

// ConnectionClosed sets this error to be returned.
func (mc *MockCount) ConnectionClosed() *Assert {
	mc.Error(ErrConnectionClosed)
	return mc.assert
}

// String representation of mocked call.
func (mc MockCount) String() string {
	return mc.assert.sprintf(`Count(ctx, "%s", %s)`, mc.argCollection, mc.argQuery)
}

// ExpectString representation of mocked call.
func (mc MockCount) ExpectString() string {
	return mc.assert.sprintf(`ExpectCount("%s", %s)`, mc.argCollection, mc.argQuery)
}
