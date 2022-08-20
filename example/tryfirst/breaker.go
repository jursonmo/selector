package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/aegis/circuitbreaker/sre"
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

	var sreOptions []sre.Option
	sreOptions = append(sreOptions, sre.WithSuccess(0.8))
	sreOptions = append(sreOptions, sre.WithRequest(1))
	sreOptions = append(sreOptions, sre.WithWindow(time.Second*5))
	sreOptions = append(sreOptions, sre.WithBucket(5))

	breaker := sre.NewBreaker(sreOptions...)
	//with breaker
	nodeSelector := firstavailable.New(selector.WithBreaker(breaker))
	nodeSelector.Apply(nodes)

	//do request
	num := 20
	for i := 0; i < num; i++ {
		//pick node
		node, done, err := nodeSelector.Select(nil)
		if err != nil {
			panic(err)
		}
		//get addr
		addr := node.Address()
		fmt.Printf("i:%d, chose addr:%s\n", i, addr)

		err = requestNode(addr)
		if done != nil {
			done(nil, selector.DoneInfo{Err: err})
		}
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func requestNode(addr string) error {
	// simulate node1 fail
	if addr == "node1_addr" {
		return errors.New("simulate node1 fail")
	}
	return nil
}
