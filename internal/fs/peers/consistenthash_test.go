package peers_test

import (
	"testing"

	"cloudborad/internal/fs"
	"cloudborad/internal/fs/peers"
)

func TestCMap(t *testing.T) {
	_ = peers.NewCMap(10, nil)
	_ = []peers.Peer{
		fs.NewDPeer("a", "http://a", 10, nil),
		fs.NewDPeer("b", "http://b", 10, nil),
	}
	// add
}
