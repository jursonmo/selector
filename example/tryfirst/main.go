package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/jursonmo/selector"
	"github.com/jursonmo/selector/firstavailable"
	"github.com/jursonmo/selector/node"
)

func main() {
	node_addrs := []string{"node1_addr", "node2_addr"}

	nodes := make([]selector.Node, 0, len(node_addrs))
	for _, addr := range node_addrs {
		nodes = append(nodes, node.New(addr, 0))
	}

	nodeSelector := firstavailable.New()
	nodeSelector.Apply(nodes)

	//do request
	num := 10
	for i := 0; i < num; i++ {
		//pick node
		node, done, err := nodeSelector.Select(nil)
		if err != nil {
			panic(err)
		}
		//get addr
		addr := node.Address()
		fmt.Printf("chose addr:%s\n", addr)

		err = requestNode(addr)
		if done != nil {
			done(nil, selector.DoneInfo{Err: err, PunishTime: 5 * time.Second})
		}
		time.Sleep(time.Second)
	}
}

func requestNode(addr string) error {
	// simulate node1 fail
	if addr == "node1_addr" {
		return errors.New("simulate node1 fail")
	}
	return nil
}
