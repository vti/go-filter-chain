package filterchain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunFilters(t *testing.T) {
	chain := FilterChain{}

	results := []int{}

	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 1)
		err := chain.Next()
		results = append(results, -1)
		return err
	}})
	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 2)
		err := chain.Next()
		results = append(results, -2)
		return err
	}})
	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 3)
		err := chain.Next()
		results = append(results, -3)
		return err
	}})

	chain.Execute()

	assert.Equal(t, 1, results[0])
	assert.Equal(t, 2, results[1])
	assert.Equal(t, 3, results[2])
	assert.Equal(t, -3, results[3])
	assert.Equal(t, -2, results[4])
	assert.Equal(t, -1, results[5])
}

func TestStopRunningOnError(t *testing.T) {
	chain := FilterChain{}

	results := []int{}

	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 1)
		err := chain.Next()
		results = append(results, -1)
		return err
	}})
	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 2)
		return errors.New("Error!")
	}})
	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 3)
		err := chain.Next()
		results = append(results, -3)
		return err
	}})

	chain.Execute()

	assert.Equal(t, 3, len(results))

	assert.Equal(t, 1, results[0])
	assert.Equal(t, 2, results[1])
	assert.Equal(t, -1, results[2])
}

func TestPropagateError(t *testing.T) {
	chain := FilterChain{}

	results := []int{}

	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 1)
		err := chain.Next()
		results = append(results, -1)
		return err
	}})
	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 2)
		return errors.New("Error!")
	}})
	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 3)
		err := chain.Next()
		results = append(results, -3)
		return err
	}})

	err := chain.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, "Error!", err.Error())
}

func TestNotRunAgain(t *testing.T) {
	chain := FilterChain{}

	results := []int{}

	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 1)
		chain.Next()
		results = append(results, -1)
		return nil
	}})

	chain.Execute()
	chain.Execute()
	chain.Execute()

	assert.Equal(t, 2, len(results))
}

func TestRewindChain(t *testing.T) {
	chain := FilterChain{}

	results := []int{}

	chain.AddFilter(&FilterFunc{func(chain *FilterChain) error {
		results = append(results, 1)
		chain.Next()
		results = append(results, -1)
		return nil
	}})

	chain.Execute()
	chain.Rewind()
	chain.Execute()
	chain.Rewind()
	chain.Execute()
	chain.Rewind()

	assert.Equal(t, 6, len(results))
}

type CustomFilter struct {
	run int
}

func (filter *CustomFilter) Execute(chain *FilterChain) error {
	filter.run++
	err := chain.Next()
	return err
}

func TestCustomStruct(t *testing.T) {
	chain := FilterChain{}

	filter := &CustomFilter{}

	chain.AddFilter(filter)

	chain.Execute()

	assert.Equal(t, 1, filter.run)
}
