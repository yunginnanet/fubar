package fubar

import "strings"

var greyStrings = []string{
	"/src/runtime",
	"/src/internal",
	"/src/testing/testing.go",
	"panic({",
	"runtime/debug",
	"testing.tRunner",
	"testing.(*T).Run",
}

var greyIfNoName = []string{
	"runtime.",
	"panic",
}

func (p *Panic) grey(s string) bool {
	for _, v := range greyStrings {
		if strings.Contains(s, v) {
			return true
		}
	}
	if strings.Contains(s, p.name) {
		return false
	}
	for _, v := range greyIfNoName {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

func splitPath(str *strings.Builder, ln string) {
	str.WriteString(reset)
	split := strings.Split(ln, "/")
	str.WriteString(brightyellow)
	one := strings.Join(split[:len(split)-1], "/")
	str.WriteString(one)
	str.WriteString("/")
	str.WriteString(reset)
	two := split[len(split)-1]

	im := 0
	if strings.Contains(two, "({") {
		im = 1
	}

	split = strings.Fields(two)

	for i, v := range split {
		if i <= im {
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
}
