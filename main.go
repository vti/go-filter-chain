// Package filtechain implements a simple FilterChain pattern.
// The filters can be either anonymous functions or custom types, they just have
// to be wrapped in something that follows Executer interface.
//
// The filter can do something before and something after calling the next
// filter, but has to propagate the return value from the next filter, or if
// needed set its own.
//
//    chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain, args ...interface{}) error {
//        // Do smth before
//        ...
//        // Call the next filter
//        err := chain.Next(args)
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
//    func (filter *CustomFilter) Execute(chain *filterchain.Chain, args ...interface{}) error {
//        fmt.Println(2)
//        err := chain.Next(args)
//        fmt.Println(-2)
//        return err
//    }
//
//    func main() {
//        chain := filterchain.New()
//
//        // Specifying filter as anon function
//        chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain, args ...interface{}) error {
//            fmt.Println(1)
//            err := chain.Next(args)
//            fmt.Println(-1)
//            return err
//        }})
//
//        // Specifying filter as a custom struct, this way it can be put in
//        // a separate package for example
//        chain.AddFilter(&CustomFilter{})
//
//        chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain, args ...interface{}) error {
//            fmt.Println(3)
//            err := chain.Next(args)
//            fmt.Println(-3)
//            return err
//        }})
//
//        chain.Execute()
//        // Will print:
//        // 1
//        // 2
//        // 3
//        // -3
//        // -2
//        // -1
//    }
package filterchain

// Executer is an interface for filters.
type Executer interface {
    Execute(*Chain, ...interface{}) error
}

// Inline is a type for adding filters as anonymous functions.
//    chain.AddFilter(&filterchain.Inline{func(chain *filterchain.Chain, args ...interface{}) error {
//        err := chain.Next(args)
//        return err
//    }})
type Inline struct {
    Handler func(*Chain, ...interface{}) error
}

// Execute runs the inlined handler.
func (filter *Inline) Execute(chain *Chain, args ...interface{}) error {
    return filter.Handler(chain, args...)
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
func (chain *Chain) Execute(args ...interface{}) error {
    pos := chain.pos
    if pos < len(chain.filters) {
        chain.pos++
        if err := chain.filters[pos].Execute(chain, args...); err != nil {
            return err
        }
    }

    return nil
}

// Next executes the next filter in the chain.
func (chain *Chain) Next(args ... interface{}) error {
    return chain.Execute(args...)
}

// Rewind rewinds the chain, so it can be run again.
func (chain *Chain) Rewind() {
    chain.pos = 0
}
