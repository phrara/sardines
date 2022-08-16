package tool

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

type PeerNode struct {
	// 节点信息
	NodeInfo *peer.AddrInfo
	// p2p节点标识
	NodeAddr multiaddr.Multiaddr
}

func (p *PeerNode) String() string {
	return p.NodeAddr.String()
}

func (p *PeerNode) ID() peer.ID {
	return p.NodeInfo.ID
}

func ParsePeerNode(p string) *PeerNode {
	if p == "" {
		return nil
	}
	addr, err := multiaddr.NewMultiaddr(p)
	if err != nil {
		return nil
	}
	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return nil
	}
	return &PeerNode{
		NodeInfo: info,
		NodeAddr: addr,
	}
}
