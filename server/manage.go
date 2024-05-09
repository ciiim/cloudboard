package server

import (
	"log"

	dlogger "cloudborad/internal/debug"
	"cloudborad/internal/dfs"
	"cloudborad/internal/dfs/peers"
)

func (s *Server) ServerInfo() (string, string) {
	return s.serverName, s._IP
}

func (s *Server) AddPeer(name, addr string) {
	s.Group.AddPeer(name, addr)
}

func (s *Server) JoinCluster(name, addr string) error {
	//boradcast to group and get all peers of the group

	dest := dfs.NewDPeerInfo(name, addr)

	//Join Cluster
	err := s.Group.PeerService.PActionTo(peers.P_ACTION_JOIN, dest)
	if err != nil {
		return err
	}

	// Get List from cluster
	peerList, err := s.Group.PeerService.GetPeerListFromPeer(dest)
	if err != nil {
		return err
	}

	//Add to peer map
	for _, peer := range peerList {
		_ = s.Group.PeerService.PHandleSyncAction(peer, peers.P_ACTION_NEW)
	}

	return nil
}

func (s *Server) QuitCluster() error {
	list := s.Group.PeerService.PList()

	return s.Group.PeerService.PActionTo(peers.P_ACTION_QUIT, list...)

}

func (s *Server) DebugOn() {
	dlogger.DebugOn()
	log.Println("[WARNING] DEBUG MODE ON")
}
