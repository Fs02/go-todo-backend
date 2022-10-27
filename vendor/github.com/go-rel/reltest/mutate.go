package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type mutate []*MockMutate

func (m *mutate) register(name string, ctxData ctxData, mutators ...rel.Mutator) *MockMutate {
	mm := &MockMutate{
		assert:      &Assert{ctxData: ctxData, repeatability: 1},
		name:        name,
		argMutators: mutators,
	}
	*m = append(*m, mm)
	return mm
}

func (m mutate) execute(name string, ctx context.Context, entity any, mutators ...rel.Mutator) error {
	for _, mm := range m {
		if (mm.argEntity == nil || reflect.DeepEqual(mm.argEntity, entity)) &&
			(mm.argEntityType == "" || mm.argEntityType == reflect.TypeOf(entity).String()) &&
			(mm.argEntityTable == "" || mm.argEntityTable == rel.NewDocument(entity, true).Table()) &&
			(mm.argEntityContains == nil || matchContains(mm.argEntityContains, entity)) &&
			(mm.argMutators == nil || matchMutators(mm.argMutators, mutators)) &&
			mm.assert.call(ctx) {
			return mm.retError
		}
	}

	mm := &MockMutate{
		assert:      &Assert{ctxData: fetchContext(ctx)},
		name:        name,
		argEntity:   entity,
		argMutators: mutators,
	}
	panic(failExecuteMessage(mm, m))
}

func (m *mutate) assert(t TestingT) bool {
	t.Helper()
	for _, mm := range *m {
		if !mm.assert.assert(t, mm) {
			return false
		}
	}

	*m = nil
	return true
}

// MockMutate asserts and simulate Insert function for test.
type MockMutate struct {
	assert            *Assert
	name              string
	argEntity         any
	argEntityType     string
	argEntityTable    string
	argEntityContains any
	argMutators       []rel.Mutator
	retError          error
}

// For assert calls for given entity.
func (mm *MockMutate) For(entity any) *MockMutate {
	mm.argEntity = entity
	return mm
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mm *MockMutate) ForType(typ string) *MockMutate {
	mm.argEntityType = "*" + strings.TrimPrefix(typ, "*")
	return mm
}

// ForTable assert calls for given table.
func (mm *MockMutate) ForTable(typ string) *MockMutate {
	mm.argEntityTable = typ
	return mm
}

// ForContains assert calls to contains some value of given struct.
func (mm *MockMutate) ForContains(contains any) *MockMutate {
	mm.argEntityContains = contains
	return mm
}

// Error sets error to be returned.
func (mm *MockMutate) Error(err error) *Assert {
	mm.retError = err
	return mm.assert
}

// Success sets no error to be returned.
func (mm *MockMutate) Success() *Assert {
	return mm.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mm *MockMutate) ConnectionClosed() *Assert {
	return mm.Error(ErrConnectionClosed)
}

// NotUnique sets not unique error to be returned.
func (mm *MockMutate) NotUnique(key string) *Assert {
	return mm.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

// String representation of mocked call.
func (mm MockMutate) String() string {
	argEntity := "<Any>"
	if mm.argEntity != nil {
		argEntity = csprint(mm.argEntity, true)
	} else if mm.argEntityContains != nil {
		argEntity = fmt.Sprintf("<Contains: %s>", csprint(mm.argEntityContains, true))
	} else if mm.argEntityType != "" {
		argEntity = fmt.Sprintf("<Type: %s>", mm.argEntityType)
	} else if mm.argEntityTable != "" {
		argEntity = fmt.Sprintf("<Table: %s>", mm.argEntityTable)
	}

	argMutators := ""
	for i := range mm.argMutators {
		argMutators += fmt.Sprintf(", %v", mm.argMutators[i])
	}

	return mm.assert.sprintf("%s(ctx, %s%s)", mm.name, argEntity, argMutators)
}

// ExpectString representation of mocked call.
func (mm MockMutate) ExpectString() string {
	argMutators := ""
	for i := range mm.argMutators {
		if i > 0 {
			argMutators += fmt.Sprintf(", %v", mm.argMutators[i])
		} else {
			argMutators += fmt.Sprintf("%v", mm.argMutators[i])
		}
	}

	return mm.assert.sprintf("Expect%s(%s).ForType(\"%T\")", mm.name, argMutators, mm.argEntity)
}
