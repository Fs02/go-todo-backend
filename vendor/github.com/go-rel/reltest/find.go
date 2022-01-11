package reltest

import (
	"context"
	"reflect"

	"github.com/go-rel/rel"
)

type find []*MockFind

func (f *find) register(ctxData ctxData, queriers ...rel.Querier) *MockFind {
	mf := &MockFind{
		assert:   &Assert{ctxData: ctxData, repeatability: 1},
		argQuery: rel.Build("", queriers...),
	}
	*f = append(*f, mf)
	return mf
}

func (f find) execute(ctx context.Context, record interface{}, queriers ...rel.Querier) error {
	query := rel.Build("", queriers...)
	for _, mf := range f {
		if matchQuery(mf.argQuery, query) &&
			mf.assert.call(ctx) {
			if mf.argRecord != nil {
				reflect.ValueOf(record).Elem().Set(reflect.ValueOf(mf.argRecord))
			}

			return mf.retError
		}
	}

	mf := &MockFind{
		assert:    &Assert{ctxData: fetchContext(ctx)},
		argQuery:  query,
		argRecord: record,
	}
	panic(failExecuteMessage(mf, f))
}

func (f *find) assert(t T) bool {
	t.Helper()
	for _, mf := range *f {
		if !mf.assert.assert(t, mf) {
			return false
		}
	}

	*f = nil
	return true
}

// MockFind asserts and simulate find function for test.
type MockFind struct {
	assert    *Assert
	argQuery  rel.Query
	argRecord interface{}
	retError  error
}

// Result sets the result of this query.
func (mf *MockFind) Result(result interface{}) *Assert {
	mf.argQuery.Table = rel.NewDocument(result, true).Table()
	mf.argRecord = result
	return mf.assert
}

// Error sets error to be returned.
func (mf *MockFind) Error(err error) *Assert {
	mf.retError = err
	return mf.assert
}

// ConnectionClosed sets this error to be returned.
func (mf *MockFind) ConnectionClosed() *Assert {
	return mf.Error(ErrConnectionClosed)
}

// NotFound sets NotFoundError to be returned.
func (mf *MockFind) NotFound() *Assert {
	return mf.Error(rel.NotFoundError{})
}

// String representation of mocked call.
func (mf MockFind) String() string {
	return mf.assert.sprintf("Find(ctx, <Any>, %s)", mf.argQuery)
}

// ExpectString representation of mocked call.
func (mf MockFind) ExpectString() string {
	return mf.assert.sprintf("ExpectFind(%s)", mf.argQuery)
}
