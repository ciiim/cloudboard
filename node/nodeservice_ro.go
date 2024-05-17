package node

import "github.com/ciiim/cloudborad/node/chash"

// ReadOnly service
type NodeServiceRO struct {
	self *Node

	//Read Only
	cMap *CMap
}

func (ns *NodeService) NodeServiceRO() *NodeServiceRO {
	ns.onceRO.Do(func() {
		if ns.ro == nil {
			ns.ro = &NodeServiceRO{
				cMap: ns.cMap,
				self: ns.self,
			}
		}
	})
	return ns.ro
}

func (ns *NodeServiceRO) Self() *Node {
	return ns.self
}

func (ns *NodeServiceRO) Pick(key []byte) *Node {
	return ns.cMap.Get(key)
}

func (ns *NodeServiceRO) GetByNodeID(nodeID string) *Node {
	return ns.cMap.GetByNodeID(nodeID)
}

func (ns *NodeServiceRO) GetAllReal() []*Node {
	decideFn := func(node chash.CHashItem) bool {
		return !node.(*chash.InnerItem).IsVirtual()
	}
	return ns.cMap.GetAll(decideFn)
}

/*
PickN

返回n个节点，包含key所属的节点以及后续的n-1个节点
*/
func (ns *NodeServiceRO) PickN(key []byte, n int) []*Node {
	return ns.cMap.GetN(key, n)
}

func (ns *NodeServiceRO) PickNext(key []byte) *Node {
	n := ns.cMap.GetN(key, 1)
	if len(n) == 2 {
		return n[1]
	}
	return n[0]
}
