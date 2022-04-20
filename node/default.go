package node

import "github.com/jursonmo/selector"

// Node is slector node
type Node struct {
	addr     string
	weight   *int64
	version  string
	name     string
	metadata map[string]string
}

// Address is node address
func (n *Node) Address() string {
	return n.addr
}

// ServiceName is node serviceName
func (n *Node) ServiceName() string {
	return n.name
}

// InitialWeight is node initialWeight
func (n *Node) InitialWeight() *int64 {
	return n.weight
}

// Version is node version
func (n *Node) Version() string {
	return n.version
}

// Metadata is node metadata
func (n *Node) Metadata() map[string]string {
	return n.metadata
}

// New node
func New(addr string, w int64) selector.Node {
	n := &Node{
		addr:   addr,
		weight: &w,
	}

	return n
}
