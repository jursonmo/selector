package random

import (
	"context"
	"math/rand"

	"github.com/jursonmo/selector"
	"github.com/jursonmo/selector/node/direct"
)

var (
	_ selector.Balancer = &Balancer{}

	// Name is balancer name
	Name = "random"
)

// Balancer is a random balancer.
type Balancer struct{}

// New random a selector.
func New() selector.Selector {
	return &selector.Default{
		Balancer:    &Balancer{},
		NodeBuilder: &direct.Builder{},
	}
}

// Pick pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	cur := rand.Intn(len(nodes))
	selected := nodes[cur]
	d := selected.Pick()
	return selected, d, nil
}
