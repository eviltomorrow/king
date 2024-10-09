package domain

type Plan2 struct {
	NotifyWithError func(error) error
	NotifyWithMsg   func(string) error

	Todo func() (string, error)

	parent *Plan2
	child  *Plan2

	signal chan struct{}
}

func (p *Plan2) Next(plan *Plan2) {
	p.child = plan
	plan.parent = p
}

func (p *Plan2) Exec() {
}

func (p *Plan2) Cancel() {
	p.signal <- struct{}{}
	for {
		child := p.child
		if child != nil {
			child.Cancel()
		}
	}
}
