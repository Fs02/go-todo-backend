package reltest

import (
	"context"
	"fmt"
	"reflect"
)

type exec []*MockExec

func (e *exec) register(ctxData ctxData, statement string, args ...any) *MockExec {
	me := &MockExec{
		assert:       &Assert{ctxData: ctxData, repeatability: 1},
		argStatement: statement,
		argArgs:      args,
	}
	*e = append(*e, me)
	return me
}

func (e exec) execute(ctx context.Context, statement string, args ...any) (int, int, error) {
	for _, me := range e {
		if me.argStatement == statement &&
			reflect.DeepEqual(me.argArgs, args) &&
			me.assert.call(ctx) {
			return me.retLastInsertedId, me.retRowsAffected, me.retError
		}
	}

	me := &MockExec{
		assert:       &Assert{ctxData: fetchContext(ctx)},
		argStatement: statement,
		argArgs:      args,
	}
	panic(failExecuteMessage(me, e))
}

func (e *exec) assert(t TestingT) bool {
	t.Helper()
	for _, me := range *e {
		if !me.assert.assert(t, me) {
			return false
		}
	}

	*e = nil
	return true
}

// MockExec asserts and simulate UpdateAny function for test.
type MockExec struct {
	assert            *Assert
	argStatement      string
	argArgs           []any
	retLastInsertedId int
	retRowsAffected   int
	retError          error
}

// Result sets the result of this query.
func (me *MockExec) Result(lastInsertedId int, rowsAffected int) *Assert {
	me.retLastInsertedId = lastInsertedId
	me.retRowsAffected = rowsAffected
	return me.assert
}

// Error sets error to be returned.
func (me *MockExec) Error(err error) *Assert {
	me.retError = err
	return me.assert
}

// ConnectionClosed sets this error to be returned.
func (me *MockExec) ConnectionClosed() *Assert {
	return me.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (me MockExec) String() string {
	args := ""
	for i := range me.argArgs {
		args += fmt.Sprintf(", %v", me.argArgs[i])
	}

	return me.assert.sprintf("Exec(ctx, \"%s\"%s)", me.argStatement, args)
}

// ExpectString representation of mocked call.
func (me MockExec) ExpectString() string {
	args := ""
	for i := range me.argArgs {
		args += fmt.Sprintf(", %v", me.argArgs[i])
	}

	return me.assert.sprintf("ExpectString(\"%s\"%s)", me.argStatement, args)
}
