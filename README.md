A simple FilterChain pattern implementation in Go.

	package main
	
	import (
	    "fmt"
	    "github.com/vti/go-filter-chain"
	)

	chain := FilterChain{}

	chain.AddFilter(&Filter{func(chain *FilterChain) error {
		fmt.Println(1)
		err := chain.Next()
		fmt.Println(-1)
		return err
	}})
	chain.AddFilter(&Filter{func(chain *FilterChain) error {
		fmt.Println(2)
		err := chain.Next()
		fmt.Println(-2)
		return err
	}})
	chain.AddFilter(&Filter{func(chain *FilterChain) error {
		fmt.Println(3)
		err := chain.Next()
		fmt.Println(-3)
		return err
	}})

	chain.Execute()

Will print:

    1
    2
    3
    -3
    -2
    -1
