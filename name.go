package fubar

import (
	"bytes"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

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

func goModName(arg0 string) string {
	d, err := os.ReadFile("./go.mod")
	if err != nil {
		return ""
	}
	if !bytes.Contains(d, module) {
		return ""
	}
	ds := bytes.Split(bytes.Split(d, module)[1], nl)
	if arg0 != "" && !strings.Contains(string(ds[0]), arg0) {
		return ""
	}
	return string(bytes.TrimPrefix(ds[0], nl))

}

func pathName() string {
	p, err := os.Getwd()
	if err != nil {
		return ""
	}
	return p
}

func myArg0() string {
	return strings.TrimSuffix(path.Base(os.Args[0]), path.Ext(os.Args[0]))

}

func pathOrGomodName() string {
	var s string

	if s = goModName(myArg0()); s != "" {
		return s
	}

	if s = myArg0(); s != "" {
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
	case strings.ReplaceAll(bi.Path, "command-line-arguments", "") != "":
		n = bi.Path
	default:
		n = pathOrGomodName()
	}
	return tail(n)
}
