package firstavailable

import (
	"context"

	"github.com/jursonmo/selector"
	"github.com/jursonmo/selector/node/available"
)

var (
	_ selector.Balancer = &Balancer{}

	// Name is balancer name, alway try the first available node
	Name = "firstavailable"
)

// Balancer is a firstavailable balancer.
type Balancer struct{}

// New firstavailable a selector.
func New() selector.Selector {
	return &selector.Default{
		Balancer:    &Balancer{},
		NodeBuilder: &available.Builder{},
	}
}

//上层resolver watch 到 地址列表后，可以 调用selector.Selector 接口 Apply 方法来设置地址, 再由NodeBuilder 把地址列表转换成自己想要的selector.WeightedNode
//到上层调用 selector.Selector Select 方法来获取一个节点时，可以调用Balancer Pick 方法来决定选择哪个WeightedNode，即Balancer Pick 方法 实现具体的算法。
//WeightedNode 负责返回一个回调函数，用于业务层反馈信息。

//比如random, 负载构建selector.Default{} 对象的主要两个部分：Balancer和NodeBuilder , NodeBuilder 用于决定构造什么样的selector.WeightedNode 实体对象
//Balancer 用于实现具体的选择算法，比如random 或者round_robin

// Pick pick a weighted node. alway try the first available node
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	var selected selector.WeightedNode
	for i := 0; i < len(nodes); i++ {
		if nodes[i].ShouldTry() {
			selected = nodes[i]
			break
		}
	}
	if selected == nil {
		selected = nodes[0]
	}
	d := selected.Pick()
	return selected, d, nil
}
