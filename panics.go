package fubar

import (
	"bufio"
	"bytes"
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
	gm     = "./go.mod"
)

type Panic struct {
	R       interface{}
	Callers []uintptr
	Trace   []byte
	name    string
}

func tail(s string) string {
	split := strings.Split(s, "/")
	switch len(split) {
	case 0:
		return ""
	case 1:
		return split[0]
	case 2:
		return split[1]
	default:
		return split[len(split)-2] + "/" + split[len(split)-1]
	}
}

func goModName() string {
	d, err := os.ReadFile("./go.mod")
	if err != nil {
		return ""
	}
	if !bytes.Contains(d, module) {
		println(string(module))
		println(string(d))
		return ""
	}
	ds := bytes.Split(bytes.Split(d, module)[1], nl)
	return string(bytes.TrimPrefix(ds[0], nl))

}

func pathName() string {
	p, err := os.Getwd()
	if err != nil {
		return ""
	}
	return p
}

func pathOrGomodName() string {
	var s string
	if s = goModName(); s != "" {
		return s
	}
	return pathName()
}

func getName() string {
	bi, _ := debug.ReadBuildInfo()
	var n string
	switch {
	case bi.Main.Path != "":
		n = bi.Main.Path
	case bi.Path != "":
		n = bi.Path
	default:
		n = pathOrGomodName()
	}
	return tail(n)
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

const (
	red          = "\033[31m"
	green        = "\033[32m"
	reset        = "\033[0m"
	gray         = "\033[90m"
	brightyellow = "\033[1;33m"
	underline    = "\033[4m"
	dim          = "\033[2m"
)

func (p *Panic) header(str *strings.Builder) {
	str.WriteString("\n")
	ln := strings.Builder{}
	ln.WriteString(gray)
	ln.WriteString(" -- ")
	ln.WriteString(reset)
	ln.WriteString(green)
	ln.WriteString(p.name)
	ln.WriteString(reset)
	ln.WriteString(" panic recovery")
	ln.WriteString(gray)
	ln.WriteString(" --")
	lnLen := ln.Len()
	space := strings.Repeat(" ", lnLen/8)
	str.WriteString(ln.String())
	ln.Reset()
	str.WriteString("\n")
	str.WriteString(space)
	str.WriteString(time.Now().Format(time.RFC3339))
	str.WriteString("\n")
	str.WriteString(strings.Repeat("-", lnLen/3+(lnLen/3)))
	str.WriteString(reset)
}

func (p *Panic) errStr(str *strings.Builder) {
	str.WriteString("\n\n")
	str.WriteString(red)
	str.WriteString(p.Error())
	str.WriteString(reset)
	str.WriteString("\n\n")
}

func head(str *strings.Builder, s string) {
	str.WriteString(gray)
	str.WriteString("----------   ")
	str.WriteString(reset)
	str.WriteString(underline)
	str.WriteString(s)
	str.WriteString(reset)
	str.WriteString("   ")
	str.WriteString(gray)
	str.WriteString("---------")
	str.WriteString(reset)
	str.WriteString("\n")

}

func (p *Panic) funcStr(str *strings.Builder) {
	head(str, "funcs")

	xerox := bufio.NewScanner(strings.NewReader(p.CallersStr()))
	for xerox.Scan() {
		str.WriteString("\n")
		if strings.Contains(xerox.Text(), p.name) && !strings.Contains(xerox.Text(), "/recovery") {
			str.WriteString(reset)
			// str.WriteString(brightyellow)
			str.WriteString(red)
			str.WriteString(underline)
			str.WriteString(xerox.Text())
			str.WriteString(reset)
			continue
		}
		if strings.Contains(xerox.Text(), "runtime.") || strings.Contains(xerox.Text(), "panic") {
			str.WriteString(reset)
			str.WriteString(dim)
			str.WriteString(gray)
			str.WriteString(xerox.Text())
			str.WriteString(reset)
			continue
		}
		str.WriteString(reset)
		str.WriteString(xerox.Text())
	}

	str.WriteString("\n\n")
}

func (p *Panic) stackStr(str *strings.Builder) {
	head(str, "stack")

	xerox := bufio.NewScanner(strings.NewReader(p.StackTraceStr()))
	for xerox.Scan() {
		str.WriteString("\n")
		if (strings.Contains(xerox.Text(), p.name) && !strings.Contains(xerox.Text(), "/recovery")) ||
			strings.Contains(xerox.Text(), "main.go") {
			str.WriteString(reset)
			split := strings.Split(xerox.Text(), "/")
			str.WriteString(brightyellow)
			one := strings.Join(split[:len(split)-1], "/")
			str.WriteString(one)
			str.WriteString("/")
			str.WriteString(reset)
			two := split[len(split)-1]

			split = strings.Fields(two)

			for i, v := range split {
				if i == 0 {
					str.WriteString(red)
					str.WriteString(underline)
					str.WriteString(v)
					str.WriteString(reset)
					continue
				}
				str.WriteString(" ")
				str.WriteString(v)
			}

			str.WriteString(reset)
			continue
		}
		if strings.Contains(xerox.Text(), p.name+"/go/pkg/mod") {
			str.WriteString(reset)
			str.WriteString(brightyellow)
			str.WriteString(xerox.Text())
			str.WriteString(reset)
			continue
		}

		if strings.Contains(xerox.Text(), "/usr/local/go/src/runtime") || strings.Contains(xerox.Text(), "panic({") ||
			strings.Contains(xerox.Text(), p.name+"/recovery") {
			str.WriteString(reset)
			str.WriteString(dim)
			str.WriteString(gray)
			str.WriteString(xerox.Text())
			str.WriteString(reset)
			continue
		}
		str.WriteString(reset)
		str.WriteString(xerox.Text())
	}

	str.WriteString("\n\n")
	str.WriteString(gray)
	str.WriteString("------------------------------")
	str.WriteString(reset)
	str.WriteString("\n")
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

func (p *Panic) delayedHardPanic() {
	go func() {
		// let our stacktrace print before we panic
		time.Sleep(100 * time.Millisecond)
		panic(p.R)
	}()
}
