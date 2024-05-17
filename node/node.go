// about node
package node

import (
	"net"

	"github.com/ciiim/cloudborad/node/chash"
)

type Node struct {
	NodeID   string `json:"node_id"`
	NodeIP   string `json:"node_ip"`
	NodePort string `json:"node_port"`
	NodeName string `json:"node_name"`
}

var _ chash.CHashItem = (*Node)(nil)

func NewNode(nodeAddr string, uniqueNodeName string) *Node {
	id := nodeAddr + uniqueNodeName
	addr, port, _ := net.SplitHostPort(nodeAddr)
	return &Node{
		NodeID:   id,
		NodeIP:   addr,
		NodePort: port,
		NodeName: uniqueNodeName,
	}
}

func (n *Node) Copy() Node {
	return *n
}

// return false if other is nil
func (n Node) Equal(other *Node) bool {
	if other == nil {
		return false
	}
	return n.ID() == other.ID()
}

func (n *Node) Compare(other chash.CHashItem) bool {
	return n.ID() == other.ID()
}

func (n Node) Name() string {
	return n.NodeName
}

func (n Node) Addr() string {
	return net.JoinHostPort(n.NodeIP, n.NodePort)
}

func (n Node) IP() string {
	return n.NodeIP
}

func (n Node) Port() string {
	return n.NodePort
}

func (n Node) ID() string {
	return n.NodeID
}
