package filterchain

type Executer interface {
    Execute(*Chain) error
}

type Func struct {
    Handler func(*Chain) error
}

func (filter *Func) Execute(chain *Chain) error {
    return filter.Handler(chain)
}

type Chain struct {
    pos int
    filters []Executer
}

func New() *Chain {
    return &Chain{}
}

func (chain *Chain) AddFilter(filter Executer) {
    chain.filters = append(chain.filters, filter)
}

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

func (chain *Chain) Next() error {
    return chain.Execute()
}

func (chain *Chain) Rewind() {
    chain.pos = 0
}
