package service

import (
	"context"
	"sardines/tool"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"io"
)

func RecvFileHandler(s network.Stream) {

	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	serv.router.AddNode(pn)
	// fmt.Println("Get a file from", pn.String())

	p := &tool.Packet{}
	header := make([]byte, tool.HEADER)
	_, err := io.ReadFull(s, header)
	if err != nil {
		return
	}

	err = p.ParseHeader(header)
	if err != nil || p.Len == 0 {
		return
	}
	val := make([]byte, p.Len)
	_, err = io.ReadFull(s, val)
	if err != nil {
		return
	}
	p.Value = val

	// fmt.Println(tool.NewFile("", ).Unwrap(p.Value).Content)

	acc := &tool.Packet{
		Tag:   1,
		Len:   3,
		Value: []byte("acc"),
	}
	data, _ := acc.Wrap()
	s.Write(data)

}

func (s *Service) SendFile(pn *tool.PeerNode, file *tool.File) bool {
	
	s.Host.Peerstore().AddAddrs(pn.ID(), pn.NodeInfo.Addrs, peerstore.PermanentAddrTTL)

	stream, err := s.Host.NewStream(context.Background(), pn.ID(), FT)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stream.Close()

	f := file.Wrap()
	packet := &tool.Packet{
		Tag:   2,
		Len:   uint32(len(f)),
		Value: f,
	}
	wrap, _ := packet.Wrap()

	stream.Write(wrap)
	

	acc := &tool.Packet{}
	header := make([]byte, tool.HEADER)
	_, err = io.ReadFull(stream, header)
	if err != nil {
		return false
	}
	err = acc.ParseHeader(header)
	if err != nil || packet.Len == 0 {
		return false
	}
	val := make([]byte, acc.Len)
	_, err = io.ReadFull(stream, val)
	if err != nil {
		return false
	}
	acc.Value = val

	if acc.ValString() == "acc" {
		fmt.Println("send a file successfully")
		return true
	} else {
		return false
	}


}
