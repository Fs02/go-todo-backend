package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type insertAll []*MockInsertAll

func (ia *insertAll) register(ctxData ctxData) *MockInsertAll {
	mia := &MockInsertAll{
		assert: &Assert{ctxData: ctxData, repeatability: 1},
	}
	*ia = append(*ia, mia)
	return mia
}

func (ia insertAll) execute(ctx context.Context, entities any) error {
	for _, mia := range ia {
		if (mia.argEntity == nil || reflect.DeepEqual(mia.argEntity, entities)) &&
			(mia.argEntityType == "" || mia.argEntityType == reflect.TypeOf(entities).String()) &&
			(mia.argEntityTable == "" || mia.argEntityTable == rel.NewCollection(entities, true).Table()) &&
			mia.assert.call(ctx) {
			return mia.retError
		}
	}

	mia := &MockInsertAll{
		assert:    &Assert{ctxData: fetchContext(ctx)},
		argEntity: entities,
	}
	panic(failExecuteMessage(mia, ia))
}

func (ia *insertAll) assert(t TestingT) bool {
	t.Helper()
	for _, mia := range *ia {
		if !mia.assert.assert(t, mia) {
			return false
		}
	}

	*ia = nil
	return true
}

// MockInsertAll asserts and simulate Insert function for test.
type MockInsertAll struct {
	assert         *Assert
	argEntity      any
	argEntityType  string
	argEntityTable string
	retError       error
}

// For assert calls for given entity.
func (mia *MockInsertAll) For(entity any) *MockInsertAll {
	mia.argEntity = entity
	return mia
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mia *MockInsertAll) ForType(typ string) *MockInsertAll {
	mia.argEntityType = "*" + strings.TrimPrefix(typ, "*")
	return mia
}

// ForTable assert calls for given table.
func (mia *MockInsertAll) ForTable(typ string) *MockInsertAll {
	mia.argEntityTable = typ
	return mia
}

// Error sets error to be returned.
func (mia *MockInsertAll) Error(err error) *Assert {
	mia.retError = err
	return mia.assert
}

// Success sets no error to be returned.
func (mia *MockInsertAll) Success() *Assert {
	return mia.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mia *MockInsertAll) ConnectionClosed() *Assert {
	return mia.Error(ErrConnectionClosed)
}

// NotUnique sets not unique error to be returned.
func (mia *MockInsertAll) NotUnique(key string) *Assert {
	return mia.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

// String representation of mocked call.
func (mia MockInsertAll) String() string {
	argEntity := "<Any>"
	if mia.argEntity != nil {
		argEntity = csprint(mia.argEntity, true)
	} else if mia.argEntityType != "" {
		argEntity = fmt.Sprintf("<Type: %s>", mia.argEntityType)
	} else if mia.argEntityTable != "" {
		argEntity = fmt.Sprintf("<Table: %s>", mia.argEntityTable)
	}

	return mia.assert.sprintf("InsertAll(ctx, %s)", argEntity)
}

// ExpectString representation of mocked call.
func (mia MockInsertAll) ExpectString() string {
	return mia.assert.sprintf("InsertAll().ForType(\"%T\")", mia.argEntity)
}
