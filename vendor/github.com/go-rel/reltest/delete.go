package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type delete []*MockDelete

func (d *delete) register(ctxData ctxData, mutators ...rel.Mutator) *MockDelete {
	md := &MockDelete{
		assert:      &Assert{ctxData: ctxData, repeatability: 1},
		argMutators: mutators,
	}
	*d = append(*d, md)
	return md
}

func (d delete) execute(ctx context.Context, entity any, mutators ...rel.Mutator) error {
	for _, md := range d {
		if (md.argEntity == nil || reflect.DeepEqual(md.argEntity, entity)) &&
			(md.argEntityType == "" || md.argEntityType == reflect.TypeOf(entity).String()) &&
			(md.argEntityTable == "" || md.argEntityTable == rel.NewDocument(entity, true).Table()) &&
			(md.argEntityContains == nil || matchContains(md.argEntityContains, entity)) &&
			matchMutators(md.argMutators, mutators) &&
			md.assert.call(ctx) {
			return md.retError
		}
	}

	md := &MockDelete{
		assert:      &Assert{ctxData: fetchContext(ctx)},
		argEntity:   entity,
		argMutators: mutators,
	}

	panic(failExecuteMessage(md, d))
}

func (d *delete) assert(t TestingT) bool {
	t.Helper()
	for _, md := range *d {
		if !md.assert.assert(t, md) {
			return false
		}
	}

	*d = nil
	return true
}

// MockDelete asserts and simulate Delete function for test.
type MockDelete struct {
	assert            *Assert
	argEntity         any
	argEntityType     string
	argEntityTable    string
	argEntityContains any
	argMutators       []rel.Mutator
	retError          error
}

// For assert calls for given entity.
func (md *MockDelete) For(entity any) *MockDelete {
	md.argEntity = entity
	return md
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (md *MockDelete) ForType(typ string) *MockDelete {
	md.argEntityType = "*" + strings.TrimPrefix(typ, "*")
	return md
}

// ForTable assert calls for given table.
func (md *MockDelete) ForTable(typ string) *MockDelete {
	md.argEntityTable = typ
	return md
}

// ForContains assert calls to contains some value of given struct.
func (md *MockDelete) ForContains(contains any) *MockDelete {
	md.argEntityContains = contains
	return md
}

// Error sets error to be returned.
func (md *MockDelete) Error(err error) *Assert {
	md.retError = err
	return md.assert
}

// Success sets no error to be returned.
func (md *MockDelete) Success() *Assert {
	return md.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (md *MockDelete) ConnectionClosed() *Assert {
	return md.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (md MockDelete) String() string {
	argEntity := "<Any>"
	if md.argEntity != nil {
		argEntity = csprint(md.argEntity, true)
	} else if md.argEntityContains != nil {
		argEntity = fmt.Sprintf("<Contains: %s>", csprint(md.argEntityContains, true))
	} else if md.argEntityType != "" {
		argEntity = fmt.Sprintf("<Type: %s>", md.argEntityType)
	} else if md.argEntityTable != "" {
		argEntity = fmt.Sprintf("<Table: %s>", md.argEntityTable)
	}

	argMutators := ""
	for i := range md.argMutators {
		argMutators += fmt.Sprintf(", %v", md.argMutators[i])
	}

	return md.assert.sprintf("Delete(ctx, %s%s)", argEntity, argMutators)
}

// ExpectString representation of mocked call.
func (md MockDelete) ExpectString() string {
	argMutators := ""
	for i := range md.argMutators {
		argMutators += fmt.Sprintf("%v", md.argMutators[i])
	}

	return md.assert.sprintf("ExpectDelete(%s).ForType(\"%T\")", argMutators, md.argEntity)
}
