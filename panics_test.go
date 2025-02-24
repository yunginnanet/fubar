package fubar

import (
	"errors"
	test_package "github.com/yunginnanet/fubar/internal/test_pkg"
	"reflect"
	"strings"
	"testing"
)

func TestNewPanic(t *testing.T) {
	testErr := errors.New("test panic")
	p := NewPanic(testErr)

	if p.R != testErr {
		t.Errorf("expected panic value %v, got %v", testErr, p.R)
	}

	if len(p.Callers) == 0 {
		t.Error("expected callers to be populated, got empty")
	}

	if len(p.Trace) == 0 {
		t.Error("expected stack trace to be populated, got empty")
	}
}

func TestPanicError(t *testing.T) {
	testErr := errors.New("test panic")
	p := NewPanic(testErr)

	if p.Error() != testErr.Error() {
		t.Errorf("expected error message '%s', got '%s'", testErr.Error(), p.Error())
	}
}

//goland:noinspection GoTypeAssertionOnErrors
func TestPanicAsError(t *testing.T) {
	testErr := errors.New("test panic")
	p := NewPanic(testErr)

	err, ok := p.AsError().(error)
	if !ok {
		t.Errorf("expected AsError to return an error, got %T", p.AsError())
	}

	if err.Error() != testErr.Error() {
		t.Errorf("expected error message '%s', got '%s'", testErr.Error(), err.Error())
	}
}

func TestPanicRecovered(t *testing.T) {
	testVal := errors.New("test value")
	p := NewPanic(testVal)

	if !reflect.DeepEqual(p.Recovered(), testVal) {
		t.Errorf("expected recovered value '%v', got '%v'", testVal, p.Recovered())
	}
}

func TestPanicStackTraceFuncStr(t *testing.T) {
	testErr := errors.New("test panic")
	p := NewPanic(testErr)

	stackTraceStr := p.CallersStr()

	if len(stackTraceStr) == 0 {
		t.Error("Expected non-empty stack trace string")
	}

	expectedFuncName := "fubar.TestPanicStackTraceFuncStr"
	if !strings.Contains(stackTraceStr, expectedFuncName) {

		t.Errorf("Expected stack trace to contain function name (may be false positive, "+
			"investigate and remove this error if so) '%s'", expectedFuncName)
	}

	t.Logf("Stack trace:\n%s", stackTraceStr)
}

func TestPanicPrint(t *testing.T) {
	testErr := errors.New("test panic")
	p := NewPanic(testErr)
	p.Stderr()
}

func TestPanicModule(t *testing.T) {
	t.Run("test_pkg", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				p := NewPanic(r)
				p.Stderr()
			}
		}()
		test_package.BadFunc()
	})
}
