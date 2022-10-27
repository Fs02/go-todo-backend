package reltest

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"github.com/go-rel/rel"
)

type iterate []*MockIterate

func (i *iterate) register(ctxData ctxData, query rel.Query, options ...rel.IteratorOption) *MockIterate {
	mi := &MockIterate{
		assert:     &Assert{ctxData: ctxData, repeatability: 1},
		argQuery:   query,
		argOptions: options,
	}
	*i = append(*i, mi)
	return mi
}

func (i iterate) execute(ctx context.Context, query rel.Query, options ...rel.IteratorOption) rel.Iterator {
	for _, mi := range i {
		if reflect.DeepEqual(mi.argOptions, options) &&
			matchQuery(mi.argQuery, query) &&
			mi.assert.call(ctx) {
			return mi
		}
	}

	mi := &MockIterate{
		assert:     &Assert{ctxData: fetchContext(ctx)},
		argQuery:   query,
		argOptions: options,
	}
	panic(failExecuteMessage(mi, i))
}

func (i *iterate) assert(t TestingT) bool {
	t.Helper()
	for _, mi := range *i {
		if !mi.assert.assert(t, mi) {
			return false
		}
	}

	*i = nil
	return true
}

type data interface {
	Len() int
	Get(index int) *rel.Document
}

// MockIterate asserts and simulate Delete function for test.
type MockIterate struct {
	assert     *Assert
	result     data
	current    int
	err        error
	argQuery   rel.Query
	argOptions []rel.IteratorOption
}

// Result sets the result of preload.
func (mi *MockIterate) Result(result any) *Assert {
	rt := reflect.TypeOf(result)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		mi.result = rel.NewCollection(result, true)
	} else {
		mi.result = rel.NewDocument(result, true)
	}
	return mi.assert
}

// Error sets error to be returned.
func (mi *MockIterate) Error(err error) *Assert {
	mi.err = err
	return mi.assert
}

// ConnectionClosed sets this error to be returned.
func (mi *MockIterate) ConnectionClosed() *Assert {
	return mi.Error(ErrConnectionClosed)
}

// Close iterator.
func (mi MockIterate) Close() error {
	return nil
}

// Next return next entity in iterator.
func (mi *MockIterate) Next(entity any) error {
	if mi.err != nil {
		return mi.err
	}

	if mi.result == nil || mi.current == mi.result.Len() {
		return io.EOF
	}

	var (
		doc = mi.result.Get(mi.current)
	)

	reflect.ValueOf(entity).Elem().Set(doc.ReflectValue())

	mi.current++
	return nil
}

// String representation of mocked call.
func (mi MockIterate) String() string {
	argOptions := ""
	for i := range mi.argOptions {
		argOptions += fmt.Sprintf(", %v", mi.argOptions[i])
	}

	return mi.assert.sprintf("Iterate(ctx, %s%s)", mi.argQuery, argOptions)
}

// ExpectString representation of mocked call.
func (mi MockIterate) ExpectString() string {
	argOptions := ""
	for i := range mi.argOptions {
		argOptions += fmt.Sprintf(", %v", mi.argOptions[i])
	}

	return mi.assert.sprintf("ExpectIterate(%s%s)", mi.argQuery, argOptions)
}
