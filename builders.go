package fubar

import (
	"bufio"
	"strings"
	"time"
)

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

	errs := p.Error()
	if strings.Contains(errs, runtimeError) {
		errs = strings.Split(p.Error(), runtimeError)[1]
		str.WriteString(gray)
		str.WriteString(runtimeError)
		str.WriteString(reset)
	}

	str.WriteString(red)
	str.WriteString(errs)
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
		if strings.Contains(xerox.Text(), p.name) {
			splitPath(str, xerox.Text())
			continue
		}
		if p.grey(xerox.Text()) {
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
		if strings.Contains(xerox.Text(), p.name) || strings.Contains(xerox.Text(), "main.go") {
			splitPath(str, xerox.Text())
			continue
		}
		if strings.Contains(xerox.Text(), p.name+"/go/pkg/mod") {
			str.WriteString(reset)
			str.WriteString(brightyellow)
			str.WriteString(xerox.Text())
			str.WriteString(reset)
			continue
		}

		if p.grey(xerox.Text()) {
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
