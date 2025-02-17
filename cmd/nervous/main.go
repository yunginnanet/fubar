// package main is a tool for testing panic recovery formatting. it's nervous.
package main

import (
	"github.com/yunginnanet/fubar"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			rec := fubar.NewPanic(r)
			rec.Stderr()
		}
	}()
	var yeet *string
	_ = *yeet
}
