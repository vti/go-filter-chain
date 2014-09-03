// Package filtechain implements a simple FilterChain pattern.
// The filters can be either anonymous functions or custom types, they just have
// to be wrapped in something that follows Executer interface.
//
// The filter can do something before and something after calling the next
// filter, but has to propagate the return value from the next filter, or if
// needed set its own.
//
//    chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
//        // Do smth before
//        ...
//        // Call the next filter
//        err := chain.Next()
//        // Do smth after
//        ...
//        // Propagate the return value
//        return err
//    }})
//
// If the current filter does not call the next filter or returns the error, the
// chain stops. This is the correct way to terminate it.
//
// Complete program example:
//
//    package main
//
//    import (
//        "fmt"
//        "github.com/vti/go-filter-chain"
//    )
//
//    type CustomFilter struct {
//    }
//
//    func (filter *CustomFilter) Execute(chain *filterchain.Chain) error {
//        fmt.Println(2)
//        err := chain.Next()
//        fmt.Println(-2)
//        return err
//    }
//
//    func main() {
//        chain := filterchain.New()
//
//        // Specifying filter as anon function
//        chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
//            fmt.Println(1)
//            err := chain.Next()
//            fmt.Println(-1)
//            return err
//        }})
//
//        // Specifying filter as a custom struct, this way it can be put in
//        // a separate package for example
//        chain.AddFilter(&CustomFilter{})
//
//        chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
//            fmt.Println(3)
//            err := chain.Next()
//            fmt.Println(-3)
//            return err
//        }})
//
//        chain.Execute()
//    }
package filterchain

// Executer is an interface for filters.
type Executer interface {
    Execute(*Chain) error
}

// Inline is a type for adding filters as anonymous functions.
//    chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain) error {
//        err := chain.Next()
//        return err
//    }})
type Inline struct {
    Handler func(*Chain) error
}

// Execute runs the inlined handler.
func (filter *Inline) Execute(chain *Chain) error {
    return filter.Handler(chain)
}

// Chain is the main type.
type Chain struct {
    pos int
    filters []Executer
}

// New creates new chain.
func New() *Chain {
    return &Chain{}
}

// AddFilter adds a filter to the chain.
func (chain *Chain) AddFilter(filter Executer) *Chain {
    chain.filters = append(chain.filters, filter)
    return chain
}

// Execute starts executing filters in the chain.
func (chain *Chain) Execute() error {
    pos := chain.pos
    if pos < len(chain.filters) {
        chain.pos++
        if err := chain.filters[pos].Execute(chain); err != nil {
            return err
        }
    }

    return nil
}

// Next executes the next filter in the chain.
func (chain *Chain) Next() error {
    return chain.Execute()
}

// Rewind rewinds the chain, so it can be run again.
func (chain *Chain) Rewind() {
    chain.pos = 0
}
