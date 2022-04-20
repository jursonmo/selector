package available

import (
	"context"
	"sync/atomic"
	"time"

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
	lastPick      int64
	available     bool
	failCnt       int //连续失败的次数
	punishTime    time.Duration
	maxPunishTime time.Duration
}

// Builder is node builder
type Builder struct {
}

// Build create node
func (b *Builder) Build(n selector.Node) selector.WeightedNode {
	return &node{Node: n, lastPick: 0, available: true, punishTime: defaultPunishT}
}

func (n *node) Pick() selector.DoneFunc {
	now := time.Now().UnixNano()
	atomic.StoreInt64(&n.lastPick, now)
	return func(ctx context.Context, di selector.DoneInfo) {
		if di.Err != nil {
			n.available = false
			if di.PunishTime != 0 {
				n.punishTime = di.PunishTime
			}
			n.failCnt++
			return
		}
		n.available = true
		n.failCnt = 0
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

func (n *node) Available() bool {
	return n.available
}

func (n *node) ShouldTry() bool {
	if n.Available() {
		return true
	}
	//惩罚时间默认是五分钟， 即失败一次，得等五分钟才能再试
	//比如超过五分钟，再 try 一次，lastPick 会被修改，还是是失败的话，failCnt++,10分钟后才会重试,以此类推下去
	//为了避免过于长时间得不到重试，设定最多60分钟会重试一次。
	elapsed := n.PickElapsed()
	if elapsed > time.Duration(n.failCnt)*n.punishTime || elapsed > 60*time.Minute {
		return true
	}
	return false
}
