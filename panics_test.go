package fubar

import (
	"bytes"
	"encoding/base64"
	"errors"
	test_package "github.com/yunginnanet/fubar/internal/test_pkg"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const exampleCode = `cGFja2FnZSBtYWluCgppbXBvcnQgImdpdGh1Yi5jb20veXVuZ2lubmFuZXQvZnViYXIiCgpmdW5jIG1haW4oKSB7CglkZWZlciBmdW5jKCkgewoJCXIgOj0gcmVjb3ZlcigpCgkJZnViYXIuSGFuZGxlUGFuaWNXaXRoRXhpdChyKQoJCS8vIG9yLCByZWNvdmVyIGFuZCBkb24ndCBleGl0OgoJCS8vIGZ1YmFyLkhhbmRsZVBhbmljKHIpCgl9KCkKCXByaW50bG4oW11zdHJpbmd7IjAiLCAiMSIsICIyIn1bM10pCn0K`

func getExampleCode() []byte {
	ba, _ := base64.StdEncoding.DecodeString(exampleCode)
	return ba
}

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

func TestHandlePanic(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		defer func() {
			r := recover()
			if HandlePanic(r) {
				t.Error("expected HandlePanic to return false")
			}
		}()
		t.Logf("yeet")
	})

	t.Run("panic", func(t *testing.T) {
		defer func() {
			r := recover()
			if !HandlePanic(r) {
				t.Error("expected HandlePanic to return true")
			}
		}()
		_ = []string{"0", "1", "2"}[3]
	})
}

func TestExampleCode(t *testing.T) {
	if gp, err := exec.LookPath("go"); gp == "" || err != nil {
		if err != nil {
			t.Logf("error looking for go in PATH: %s", err.Error())
		}
		t.Skip("go not found in PATH, skipping test")
	}

	tdir := t.TempDir()
	tfile := filepath.Join(tdir, "main.go")

	t.Run("exit_1", func(t *testing.T) {
		t.Logf("writing example code to %s:\n%s\n", tfile, getExampleCode())

		if err := os.WriteFile(tfile, getExampleCode(), 0644); err != nil {
			t.Fatalf("error writing example code to file: %s", err.Error())
		}

		var err error

		var cmd = exec.Command("go", "run", tfile)
		t.Logf("$ %s", cmd.String())
		if err = cmd.Run(); err == nil {
			t.Error("expected example code to panic and return an error")
		}

		if !strings.Contains(err.Error(), "exit status 1") {
			t.Error("expected example code to exit with status 1")
		}

		t.Logf("example code exited with error: %s", err.Error())
	})

	t.Run("exit_0", func(t *testing.T) {
		code := bytes.ReplaceAll(getExampleCode(), []byte("}[3])"), []byte("}[2])"))

		t.Logf("writing example code to %s:\n%s\n", tfile, code)

		if err := os.WriteFile(tfile, code, 0644); err != nil {
			t.Fatalf("error writing example code to file: %s", err.Error())
		}

		var err error

		var cmd = exec.Command("go", "run", tfile)
		t.Logf("$ %s", cmd.String())
		if err = cmd.Run(); err != nil {
			t.Errorf("expected moddified example code to run successfully, got error: %s", err.Error())
		}

		t.Logf("modified example code exited with success")
	})

	t.Run("exit_1_has_name", func(t *testing.T) {
		t.Logf("writing example code to %s:\n%s\n", tfile, getExampleCode())

		if err := os.WriteFile(tfile, getExampleCode(), 0644); err != nil {
			t.Fatalf("error writing example code to file: %s", err.Error())
		}

		var err error

		buildPath := filepath.Join(tdir, "yeeterson_mcgee")

		var cmd = exec.Command("go", "build", "-o", buildPath, "-trimpath", tfile)
		t.Logf("$ %s", cmd.String())
		if err = cmd.Run(); err != nil {
			t.Skipf("expected example code to build successfully, got error: %s", err.Error())
		}

		var output []byte

		cmd = exec.Command(buildPath)
		if output, err = cmd.CombinedOutput(); err == nil {
			t.Error("expected example code to panic and return an error")
		}

		t.Logf("output:\n\n\n\n%s\n\n\n\n", string(output))

		if !strings.Contains(err.Error(), "exit status 1") {
			t.Error("expected example code to exit with status 1")
		}

		if !strings.Contains(string(output), "yeeterson_mcgee") {
			t.Error("expected output to contain name of executable")
		}
	})
}
