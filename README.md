# fubar

pretty panic recovery for golang.

## Usage

```go
package main

import "github.com/yunginnanet/fubar"

func main() {
	defer func() {
		r := recover()
		fubar.HandlePanic(r)
		// or, exit status 1 after printing:
		// fubar.HandlePanicWithExit(r)
	}()
	println([]string{"0", "1", "2"}[3])
}

```

### Output

![call5](assets/fubar.gif)

## GoDoc

```go
const NumCallersCaught = 48
```

#### func  HandlePanic

```go
func HandlePanic(r interface{}) bool
```

#### func  HandlePanicWithExit

```go
func HandlePanicWithExit(r interface{}) bool
```

#### type Panic

```go
type Panic struct {
	R       interface{}
	Callers []uintptr
	Trace   []byte
}
```


#### func  NewPanic

```go
func NewPanic(r interface{}) *Panic
```

#### func (*Panic) AsError

```go
func (p *Panic) AsError() error
```

#### func (*Panic) CallersStr

```go
func (p *Panic) CallersStr() string
```

#### func (*Panic) Error

```go
func (p *Panic) Error() string
```

#### func (*Panic) Recovered

```go
func (p *Panic) Recovered() interface{}
```

#### func (*Panic) StackTraceStr

```go
func (p *Panic) StackTraceStr() string
```

#### func (*Panic) Stderr

```go
func (p *Panic) Stderr()
```

#### func (*Panic) Stdout

```go
func (p *Panic) Stdout()
```

#### func (*Panic) String

```go
func (p *Panic) String() string
```
