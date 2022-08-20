package circuitBreaker

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/jursonmo/selector"
)

const (
	defaultWeight  = 100
	defaultPunishT = 5 * time.Minute
)

var (
	_ selector.WeightedNode        = &node{}
	_ selector.WeightedNodeBuilder = &Builder{}
)

// node is endpoint instance
type node struct {
	selector.Node

	// last lastPick timestamp
	lastPick int64
	breaker  circuitbreaker.CircuitBreaker
}

// Builder is node builder
type Builder struct {
	breaker circuitbreaker.CircuitBreaker
}

func NewBuilder(b circuitbreaker.CircuitBreaker) *Builder {
	return &Builder{breaker: b}
}

// Build create node
func (b *Builder) Build(n selector.Node) selector.WeightedNode {
	return &node{Node: n, lastPick: 0, breaker: b.breaker}
}

func (n *node) Pick() selector.DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	return func(ctx context.Context, di selector.DoneInfo) {
		//通过回调来设置breaker 失败或成功的次数。
		//所以业务层不管结果如何，必须调用DoneFunc回调。
		if di.Err != nil {
			n.breaker.MarkFailed()
			return
		}
		n.breaker.MarkSuccess()
	}
}

// Weight is node effective weight
func (n *node) Weight() float64 {
	if n.InitialWeight() != nil {
		return float64(*n.InitialWeight())
	}
	return defaultWeight
}

func (n *node) PickElapsed() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&n.lastPick))
}

// if breaker don't allow, return false, and should check the next node
func (n *node) ShouldTry() bool {
	if err := n.breaker.Allow(); err != nil {
		return false
	}
	return true
}
