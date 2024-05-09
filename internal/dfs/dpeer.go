package dfs

import (
	"context"
	"hash/crc64"
	"log"
	"math/rand"
	"strings"
	"time"

	dlogger "cloudborad/internal/debug"
	"cloudborad/internal/dfs/peers"
	"cloudborad/internal/dfs/peers/chash"
)

type DAddr string

func (a DAddr) String() string {
	return string(a)
}

func (a DAddr) Port() string {
	t := strings.Split(string(a), ":")
	if len(t) != 2 {
		return ""
	}
	return t[len(t)-1]
}

func (a DAddr) IP() string {
	t := strings.Split(string(a), ":")
	if len(t) != 2 {
		return ""
	}
	return t[0]
}

type DPeer struct {
	info         *DPeerInfo
	hashMap      *chash.CMap
	syncSettings SyncSettings
	syncRand     *rand.Rand
}

var _ peers.Peer = (*DPeer)(nil)

type DPeerInfo struct {
	PeerID       int64              `json:"peer_id"`
	PeerName     string             `json:"peer_name"`
	PeerAddr     peers.Addr         `json:"peer_addr"` //include port e.g. 10.10.1.5:9631
	PeerStat     peers.PeerStatType `json:"peer_stat"`
	LastPingTime time.Time
	Version      int64
}

func NewDPeerInfo(name, addr string) *DPeerInfo {
	return &DPeerInfo{
		PeerID:   int64(crc64.Checksum([]byte(name+time.Now().String()), crc64.MakeTable(crc64.ECMA))),
		PeerName: name,
		PeerAddr: DAddr(addr),
		PeerStat: peers.P_STAT_ONLINE,
		Version:  0,
	}
}

var _ peers.PeerInfo = (*DPeerInfo)(nil)

func NewDPeer(name, addr string, replicas int, fn chash.CHash, settings SyncSettings) *DPeer {
	dlogger.Dlog.LogDebugf("NewDPeer", "name: %s, addr: %s", name, addr)
	p := &DPeer{
		info:         NewDPeerInfo(name, addr),
		hashMap:      chash.NewCMap(replicas, fn),
		syncSettings: settings,
	}

	//Add self to hashMap
	p.hashMap.Add(*p.info)
	return p
}

func (p *DPeer) PName() string {
	return p.info.PeerName
}

func (p *DPeer) PAddr() peers.Addr {
	return p.info.PeerAddr
}

func (p *DPeer) Pick(key string) peers.PeerInfo {
	return p.hashMap.Get(key)
}

func (p *DPeer) PAdd(pis ...peers.PeerInfo) {
	p.hashMap.Add(pis...)
}

func (p *DPeer) PDel(pis ...peers.PeerInfo) {
	p.hashMap.Del(pis...)
}

/*
recieve peer action from other peer
source peer - pi_in
*/
func (p *DPeer) PHandleSyncAction(pi_in peers.PeerInfo, action peers.PeerActionType) error {
	dlogger.Dlog.LogDebugf("PSync", "pi_in: %v, action: %s", pi_in, action.String())
	if pi_in.Equal(p.info) {
		log.Println("[DPeer] Cannot Operate myself")
		return nil
	}
	var err error
	switch action {
	case peers.P_ACTION_JOIN:
		// notify other peers - action P_ACTION_NEW
		client := newRPCPeerClient()
		ctx, cancel := context.WithTimeout(context.Background(), _RPC_TIMEOUT)
		defer cancel()
		list := p.PList()
		err = client.peerActionTo(ctx, pi_in, peers.P_ACTION_NEW, list...)

	case peers.P_ACTION_QUIT:
		// remove peer from hashMap
		p.PDel(pi_in)

	case peers.P_ACTION_NEW:
		// add peer to hashMap
		p.PAdd(pi_in)

	default:
		log.Println("[DPeer] Unknown action")
	}

	return err
}

/*
send peer action to other peer

pi_to - destination peer
*/
func (p *DPeer) PActionTo(action peers.PeerActionType, pi_to ...peers.PeerInfo) error {
	dlogger.Dlog.LogDebugf("PActionTo", "action: %s, pi_to: %v", action.String(), pi_to)
	client := newRPCPeerClient()
	ctx, cancel := context.WithTimeout(context.Background(), _RPC_TIMEOUT)
	defer cancel()
	return client.peerActionTo(ctx, p.info, action, pi_to...)
}

func (p *DPeer) GetPeerListFromPeer(pi peers.PeerInfo) ([]peers.PeerInfo, error) {
	client := newRPCPeerClient()
	ctx, cancel := context.WithTimeout(context.Background(), _RPC_TIMEOUT)
	defer cancel()
	list, err := client.getPeerList(ctx, pi)
	if err != nil {
		return nil, err
	}
	peerList := make([]peers.PeerInfo, 0, len(list))
	for _, v := range list {
		peerList = append(peerList, NewDPeerInfo(v.PName(), v.PAddr().String()))
	}
	return peerList, nil
}

func (p *DPeer) PNext(key string) peers.PeerInfo {
	return p.hashMap.GetPeerNext(key, 1)
}

func (p *DPeer) PList() []peers.PeerInfo {
	return p.hashMap.List()
}

func (p *DPeer) Info() peers.PeerInfo {
	return *p.info
}

func (pi DPeerInfo) Equal(other peers.PeerInfo) bool {
	o, ok := other.(DPeerInfo)
	if !ok {
		o, ok := other.(*DPeerInfo)
		if !ok {
			return false
		}
		return pi.PeerName == o.PeerName && pi.PeerAddr == o.PeerAddr
	}
	return pi.PeerName == o.PeerName && pi.PeerAddr == o.PeerAddr
}

func (pi DPeerInfo) PID() int64 {
	return pi.PeerID
}

func (pi DPeerInfo) PName() string {
	return pi.PeerName
}

func (pi DPeerInfo) PAddr() peers.Addr {
	return pi.PeerAddr
}

func (pi DPeerInfo) PVersion() int64 {
	return pi.Version
}

func (pi DPeerInfo) PStat() peers.PeerStatType {
	return pi.PeerStat
}

func PeerInfoListToDpeerInfoList(list []peers.PeerInfo) []DPeerInfo {
	dpeerList := make([]DPeerInfo, 0, len(list))
	for _, v := range list {
		dpeerList = append(dpeerList, v.(DPeerInfo))
	}
	return dpeerList
}
