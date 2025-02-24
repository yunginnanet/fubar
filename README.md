# FUBAR

Pretty panic recovery for golang.

### Usage

```go
package main

import "github.com/yunginnanet/fubar"

func main() {
	defer func() {
		if r := recover(); r != nil {
			p := fubar.NewPanic(r)
			p.Stderr()
		}
	}()

	strs := []string{"0", "1", "2"}
	println(strs[3])
}
```

### GIF

![call5](assets/fubar.gif)
