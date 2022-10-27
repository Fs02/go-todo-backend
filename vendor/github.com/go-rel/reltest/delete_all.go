package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type deleteAll []*MockDeleteAll

func (da *deleteAll) register(ctxData ctxData) *MockDeleteAll {
	mda := &MockDeleteAll{
		assert: &Assert{ctxData: ctxData, repeatability: 1},
	}
	*da = append(*da, mda)
	return mda
}

func (da deleteAll) execute(ctx context.Context, entity any) error {
	for _, mda := range da {
		if (mda.argEntity == nil || reflect.DeepEqual(mda.argEntity, entity)) &&
			(mda.argEntityType == "" || mda.argEntityType == reflect.TypeOf(entity).String()) &&
			(mda.argEntityTable == "" || mda.argEntityTable == rel.NewCollection(entity, true).Table()) &&
			mda.assert.call(ctx) {
			return mda.retError
		}
	}

	mda := &MockDeleteAll{
		assert:    &Assert{ctxData: fetchContext(ctx)},
		argEntity: entity,
	}
	panic(failExecuteMessage(mda, da))
}

func (da *deleteAll) assert(t TestingT) bool {
	t.Helper()
	for _, mda := range *da {
		if !mda.assert.assert(t, mda) {
			return false
		}
	}

	*da = nil
	return true
}

// MockDeleteAll asserts and simulate Delete function for test.
type MockDeleteAll struct {
	assert         *Assert
	argEntity      any
	argEntityType  string
	argEntityTable string
	retError       error
}

// For assert calls for given entity.
func (mda *MockDeleteAll) For(entity any) *MockDeleteAll {
	mda.argEntity = entity
	return mda
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mda *MockDeleteAll) ForType(typ string) *MockDeleteAll {
	mda.argEntityType = "*" + strings.TrimPrefix(typ, "*")
	return mda
}

// ForTable assert calls for given table.
func (mda *MockDeleteAll) ForTable(typ string) *MockDeleteAll {
	mda.argEntityTable = typ
	return mda
}

// Error sets error to be returned.
func (mda *MockDeleteAll) Error(err error) *Assert {
	mda.retError = err
	return mda.assert
}

// Success sets no error to be returned.
func (mda *MockDeleteAll) Success() *Assert {
	return mda.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mda *MockDeleteAll) ConnectionClosed() *Assert {
	return mda.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (mda MockDeleteAll) String() string {
	argEntity := "<Any>"
	if mda.argEntity != nil {
		argEntity = csprint(mda.argEntity, true)
	} else if mda.argEntityType != "" {
		argEntity = fmt.Sprintf("<Type: %s>", mda.argEntityType)
	} else if mda.argEntityTable != "" {
		argEntity = fmt.Sprintf("<Table: %s>", mda.argEntityTable)
	}

	return mda.assert.sprintf("DeleteAll(ctx, %s)", argEntity)
}

// ExpectString representation of mocked call.
func (mda MockDeleteAll) ExpectString() string {
	return mda.assert.sprintf(`ExpectDeleteAll().ForType("%T")`, mda.argEntity)
}
