
看了 [kratos selector](https://github.com/go-kratos/kratos/tree/v2.2.1/selector) 源码

##### 大概的流程是： 
上层resolver watch 到 地址列表后，可以 调用selector.Selector 接口 Apply 方法来设置地址, 再由NodeBuilder 把地址列表转换成自己想要的selector.WeightedNode
上层通过调用 selector.Selector 的 Select方法来尝试获取一个节点时，就会调用Balancer Pick 方法来决定选择哪个WeightedNode，即Balancer Pick 方法 实现具体的算法。

WeightedNode 负责返回一个回调函数DoneFunc func(ctx context.Context, di DoneInfo)，用于业务层反馈信息DoneInfo 给node, 从而可以根据反馈的信息DoneInfo 来更新 node 状态。

比如random, 负载构建selector.Default{} 对象的主要两个部分：Balancer和NodeBuilder , NodeBuilder 用于决定构造什么样的selector.WeightedNode 实体对象
Balancer 用于实现具体的随机选择算法。当然 如果要实现round_robin, 也是同样的逻辑流程.

##### 项目中
实际的项目中，有两个节点服务器，比如 华北节点，华南节点， 广东地区的尽量上报给华南节点，如果华南节点有异常才考虑上报给华北节点，即对于需要上报数据的机器来说，这两个节点是主备关系。
对于上报的机器来说，尽量尝试上报给第一个节点。在 [kratos selector](github.com/go-kratos/kratos/v2/selector) 基础上, 我自己添加一种“alway try the first available node” , 如果第一个节点异常后，业务层可以通过DoneInfo 的punishTime 来给这个失败的节点一个惩罚时间，即在这个惩罚时间内，不会再上报数据给这个节点。[example](https://github.com/jursonmo/selector/tree/master/example/tryfirst)

##### TODO
当节点异常后时，给节点一个惩罚时间直接熔断，但这个时间不好定， 可以节点的错误率到达一定值才熔断，可以借鉴 google SRE 过载保护算法实现熔断器：
rejectProba = max(0,(requests−K∗accepts)/(requests+1))



