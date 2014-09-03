package filterchain

type Filter interface {
    Execute(*FilterChain) error
}

type FilterFunc struct {
    handler func(*FilterChain) error
}

func (filter *FilterFunc) Execute(chain *FilterChain) error {
    return filter.handler(chain)
}

type FilterChain struct {
    pos int
    filters []Filter
}

func (chain *FilterChain) AddFilter(filter Filter) {
    chain.filters = append(chain.filters, filter)
}

func (chain *FilterChain) Execute() error {
    pos := chain.pos
    if pos < len(chain.filters) {
        chain.pos++
        if err := chain.filters[pos].Execute(chain); err != nil {
            return err
        }
    }

    return nil
}

func (chain *FilterChain) Next() error {
    return chain.Execute()
}

func (chain *FilterChain) Rewind() {
    chain.pos = 0
}
