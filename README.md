# filterchain
--
    import "github.com/vti/go-filter-chain"

Package filtechain implements a simple FilterChain pattern. The filters can be
either anonymous functions or custom types, they just have to be wrapped in
something that follows Executer interface.

The filter can do something before and something after calling the next filter,
but has to propagate the return value from the next filter, or if needed set its
own.

    chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
        // Do smth before
        ...
        // Call the next filter
        err := chain.Next()
        // Do smth after
        ...
        // Propagate the return value
        return err
    }})

If the current filter does not call the next filter or returns the error, the
chain stops. This is the correct way to terminate it.

Complete program example:

    package main

    import (
        "fmt"
        "github.com/vti/go-filter-chain"
    )

    type CustomFilter struct {
    }

    func (filter *CustomFilter) Execute(chain *filterchain.Chain) error {
        fmt.Println(2)
        err := chain.Next()
        fmt.Println(-2)
        return err
    }

    func main() {
        chain := filterchain.New()

        // Specifying filter as anon function
        chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
            fmt.Println(1)
            err := chain.Next()
            fmt.Println(-1)
            return err
        }})

        // Specifying filter as a custom struct, this way it can be put in
        // a separate package for example
        chain.AddFilter(&CustomFilter{})

        chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
            fmt.Println(3)
            err := chain.Next()
            fmt.Println(-3)
            return err
        }})

        chain.Execute()
    }

## Usage

#### type Chain

```go
type Chain struct {
}
```

Chain is the main type.

#### func  New

```go
func New() *Chain
```
New creates new chain.

#### func (*Chain) AddFilter

```go
func (chain *Chain) AddFilter(filter Executer) *Chain
```
AddFilter adds a filter to the chain.

#### func (*Chain) Execute

```go
func (chain *Chain) Execute() error
```
Execute starts executing filters in the chain.

#### func (*Chain) Next

```go
func (chain *Chain) Next() error
```
Next executes the next filter in the chain.

#### func (*Chain) Rewind

```go
func (chain *Chain) Rewind()
```
Rewind rewinds the chain, so it can be run again.

#### type Executer

```go
type Executer interface {
	Execute(*Chain) error
}
```

Executer is an interface for filters.

#### type Inline

```go
type Inline struct {
	Handler func(*Chain) error
}
```

Inline is a type for adding filters as anonymous functions.

    chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
        err := chain.Next()
        return err
    }})

#### func (*Inline) Execute

```go
func (filter *Inline) Execute(chain *Chain) error
```
Execute runs the inlined handler.
