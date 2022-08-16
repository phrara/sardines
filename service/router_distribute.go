package service

import (
	"context"
	"sardines/tool"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"io"
	"sync"

	"github.com/libp2p/go-libp2p-core/network"
)

// Router Distribution
// Router table will be distributed periodically
// When peers get the distributed router tables, they use then to update their own router tables
// This automatically renew the router info of the decentralized network

func RouterDistributeHandler(s network.Stream) {
	pn := tool.ParsePeerNode(s.Conn().RemoteMultiaddr().String() + "/p2p/" + s.Conn().RemotePeer().String())
	serv.router.AddNode(pn)
	fmt.Println("Get a distributed router table from", pn.String())

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

	defer s.Close()
	if p.ValString() == "\n" {
		return
	} else {
		//Parse the raw data and use it to update the local router table
		fmt.Println("remote router info is: \n", p.ValString())
		data := serv.router.ParseData(p.ValString())
		serv.router.Update(data)
	}
}

func (s *Service) RouterDistribute() (errNum int) {

	var wg sync.WaitGroup

	localRouter := serv.router.RawData()
	fmt.Println("recent router table:\n", string(localRouter))
	nodes := serv.router.AllNodes()
	for e := nodes.Front(); e != nil; e = e.Next() {
		pn := e.Value.(*tool.PeerNode)

		s.Host.Peerstore().AddAddrs(pn.ID(), pn.NodeInfo.Addrs, peerstore.PermanentAddrTTL)

		wg.Add(1)
		go func(p *tool.PeerNode) {
			defer wg.Done()

			stream, err := s.Host.NewStream(context.Background(), p.ID(), RD)
			if err != nil {
				fmt.Println(err)
				errNum = errNum + 1
				return
			}

			packet := &tool.Packet{
				Tag:   2,
				Len:   uint32(len(localRouter)),
				Value: localRouter,
			}
			wrap, _ := packet.Wrap()

			stream.Write(wrap)

		}(pn)
	}

	wg.Wait()
	return errNum
}
