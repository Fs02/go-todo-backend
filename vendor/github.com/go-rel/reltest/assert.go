package reltest

import (
	"context"
	"fmt"
	"reflect"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Logf(format string, args ...any)
	Errorf(format string, args ...any)
	Helper()
}

type Assert struct {
	ctxData       ctxData
	repeatability int // 0 means not limited
	totalCalls    int
	optional      bool
}

// Once set max calls to one time.
func (a *Assert) Once() {
	a.Times(1)
}

// Twice set max calls to two times.
func (a *Assert) Twice() {
	a.Times(2)
}

// Many set max calls to unlimited times.
func (a *Assert) Many() {
	a.Times(0)
}

// Times set number of allowed calls.
func (a *Assert) Times(times int) {
	a.repeatability = times
}

// Maybe allow calls to be skipped.
func (a *Assert) Maybe() {
	a.optional = true
}

// this function needs to be called as last condition of if
// otherwise recorded total calls will be wrong
func (a *Assert) call(ctx context.Context) bool {
	if a.ctxData != fetchContext(ctx) || (a.repeatability != 0 && a.totalCalls >= a.repeatability) {
		return false
	}

	a.totalCalls++
	return true
}

func (a Assert) assert(t TestingT, mock any) bool {
	if a.optional ||
		(a.repeatability == 0 && a.totalCalls > 0) ||
		(a.repeatability != 0 && a.totalCalls >= a.repeatability) {
		return true
	}

	t.Helper()
	if a.repeatability > 1 {
		t.Errorf("FAIL: Need to make %d more call(s) to satisfy mock:\n%s", a.repeatability-a.totalCalls, mock)
	} else {
		t.Errorf("FAIL: Mock defined but not called:\n%s", mock)
	}

	return false
}

func (a Assert) sprintf(format string, args ...any) string {
	if a.ctxData.txDepth != 0 {
		return a.ctxData.String() + " " + fmt.Sprintf(format, args...)
	}
	return fmt.Sprintf(format, args...)
}

func failExecuteMessage(call any, mocks any) string {
	var (
		mocksStr      string
		callStr       = call.(interface{ String() string }).String()
		expectCallStr = call.(interface{ ExpectString() string }).ExpectString()
		rv            = reflect.ValueOf(mocks)
	)

	for i := 0; i < rv.Len(); i++ {
		mocksStr += fmt.Sprintf("\n\t- %s", rv.Index(i).Interface())
	}

	if mocksStr == "" {
		mocksStr = "None"
	}

	return fmt.Sprintf("FAIL: this call is not mocked:\n\t%s\nMaybe try adding mock:\n\t%s\n\nMocked calls:%s\n\n", callStr, expectCallStr, mocksStr)
}
