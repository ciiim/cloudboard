package node

import (
	"github.com/ciiim/cloudborad/node/chash"
)

// 并发安全
type CMap struct {
	*chash.ConsistentHash
}

// create a new consistent hash map
func NewCMap(replicas int, fn chash.ConsistentHashFn) *CMap {
	return &CMap{
		chash.NewConsistentHash(replicas, fn),
	}
}

func (c *CMap) GetByNodeID(nodeID string) *Node {
	node := chash.NewQuerier[*Node](c.ConsistentHash).GetByID(nodeID)
	if node.ID() == "nil" {
		return nil
	}
	return node
}

func (c *CMap) Get(key []byte) *Node {
	node := chash.NewQuerier[*Node](c.ConsistentHash).Get(key)
	if node.ID() == "nil" {
		return nil
	}
	return node
}

func (c *CMap) GetN(key []byte, n int) []*Node {
	nodes := chash.NewQuerier[*Node](c.ConsistentHash).GetN(key, n)
	if len(nodes) == 0 || nodes[0].ID() == "nil" {
		return []*Node{}
	}
	return nodes
}

func (c *CMap) GetAll(decideFn func(chash.CHashItem) bool) []*Node {
	items := chash.NewQuerier[*Node](c.ConsistentHash).GetAll(decideFn)
	nodes := make([]*Node, 0, len(items))
	for _, item := range items {
		nodes = append(nodes, item)
	}
	return nodes
}
