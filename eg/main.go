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
