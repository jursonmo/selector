package selector

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/aegis/circuitbreaker"
)

// SelectOptions is Select Options.
type SelectOptions struct {
	Filters []Filter
	Breaker circuitbreaker.CircuitBreaker
}

// SelectOption is Selector option.
type SelectOption func(*SelectOptions)

// Filter is node filter function.
type Filter func(context.Context, []Node) []Node

// WithFilter with filter options
func WithFilter(fn ...Filter) SelectOption {
	return func(opts *SelectOptions) {
		opts.Filters = fn
	}
}

func WithBreaker(b circuitbreaker.CircuitBreaker) SelectOption {
	return func(o *SelectOptions) {
		o.Breaker = b
	}
}

// ErrNoAvailable is no available node.
var ErrNoAvailable = errors.New("no_available_node")

// Selector is node pick balancer.
type Selector interface {
	Rebalancer

	// Select nodes
	// if err == nil, selected and done must not be empty.
	Select(ctx context.Context, opts ...SelectOption) (selected Node, done DoneFunc, err error)
}

// Rebalancer is nodes rebalancer.
type Rebalancer interface {
	// apply all nodes when any changes happen
	Apply(nodes []Node)
}

// Node is node interface.
type Node interface {
	// Address is the unique address under the same service
	Address() string

	// ServiceName is service name
	ServiceName() string

	// InitialWeight is the initial value of scheduling weight
	// if not set return nil
	InitialWeight() *int64

	// Version is service node version
	Version() string

	// Metadata is the kv pair metadata associated with the service instance.
	// version,namespace,region,protocol etc..
	Metadata() map[string]string
}

// DoneInfo is callback info when RPC invoke done.
type DoneInfo struct {
	// Response Error
	Err error

	//todo
	PunishTime time.Duration
}

// ReplyMeta is Reply Metadata.
type ReplyMeta interface {
	Get(key string) string
}

// DoneFunc is callback function when RPC invoke done.
type DoneFunc func(ctx context.Context, di DoneInfo)
