package fubar

import (
	"errors"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const NumCallersCaught = 48

var (
	nl     = []byte{0x0a}
	module = []byte{0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x20}
)

const runtimeError = "runtime error: "

type Panic struct {
	R       interface{}
	Callers []uintptr
	Trace   []byte
	name    string
}

func NewPanic(r interface{}) *Panic {
	var callers [NumCallersCaught]uintptr
	found := runtime.Callers(0, callers[:])
	if rStr, strOk := r.(string); strOk {
		r = errors.New(rStr)
	}

	return &Panic{
		R:       r,
		Callers: callers[:found],
		Trace:   debug.Stack(),
		name:    getName(),
	}
}

func (p *Panic) Error() string {
	return p.R.(error).Error()
}

func (p *Panic) AsError() error {
	return p.R.(error)
}

func (p *Panic) Recovered() interface{} {
	return p.R
}

func (p *Panic) CallersStr() string {
	str := &strings.Builder{} // lol lets not use a pool for panic recovery
	frames := runtime.CallersFrames(p.Callers)
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		str.WriteString(frame.Function)
		str.WriteString("\n")
	}
	return str.String()
}

func (p *Panic) StackTraceStr() string {
	return string(p.Trace)
}

func (p *Panic) String() string {
	if os.Getenv("FUBAR_HARD_PANIC") != "" {
		defer p.delayedHardPanic()
	}

	str := &strings.Builder{}

	p.header(str)

	p.errStr(str)

	p.funcStr(str)

	p.stackStr(str)

	return str.String()
}

func (p *Panic) Stderr() {
	_, _ = os.Stderr.WriteString(p.String())
}

func (p *Panic) Stdout() {
	_, _ = os.Stdout.WriteString(p.String())
}

func HandlePanic(r interface{}) bool {
	if r == nil {
		return false
	}
	p := NewPanic(r)
	p.Stderr()
	return true
}

//goland:noinspection GoUnusedExportedFunction (unit test compiles exe that uses and tests function)
func HandlePanicWithExit(r interface{}) bool {
	if r == nil {
		return false
	}
	p := NewPanic(r)
	p.Stderr()
	os.Exit(1)
	return true // unreachable
}

func (p *Panic) delayedHardPanic() {
	go func() {
		// let our stacktrace print before we panic
		time.Sleep(100 * time.Millisecond)
		panic(p.R)
	}()
}
