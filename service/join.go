package service

import (
	"context"
	"sardines/tool"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"io"
)

func JoinApplyHandler(s network.Stream) {
	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	fmt.Println("receive a join application from", pn.String())

	// 节点加入路由表
	serv.router.AddNode(pn)

	p := &tool.Packet{
		Tag:   1,
		Len:   3,
		Value: []byte("acc"),
	}
	data, _ := p.Wrap()
	s.Write(data)

}

func (s *Service) JoinApply(bootstrapNode *tool.PeerNode) bool {
	s.Host.Peerstore().AddAddrs(bootstrapNode.ID(), bootstrapNode.NodeInfo.Addrs, peerstore.PermanentAddrTTL)

	stream, err := s.Host.NewStream(context.Background(), bootstrapNode.ID(), JOIN)
	if err != nil {
		fmt.Println(err)
		return false
	}
	// parse the header
	packet := &tool.Packet{}
	header := make([]byte, tool.HEADER)
	_, err = io.ReadFull(stream, header)
	if err != nil {
		return false
	}
	err = packet.ParseHeader(header)
	if err != nil || packet.Len == 0 {
		return false
	}
	// get body
	val := make([]byte, packet.Len)
	_, err = io.ReadFull(stream, val)
	if err != nil {
		return false
	}
	packet.Value = val

	defer stream.Close()

	if packet.ValString() == "acc" {
		return true
	} else {
		return false
	}
}
