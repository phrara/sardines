package router

import (
	"bytes"
	"container/list"
	"github.com/libp2p/go-libp2p/core/peer"
	"sardines/tool"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

var rt *Router

func init() {
	rt = &Router{}
}

type Table map[int]*list.List

type Router struct {
	host host.Host
}

func New(hn host.Host) *Router {
	rt.host = hn
	return rt
}

func (r *Router) AddNode(pn *tool.PeerNode) {
	r.host.Peerstore().AddAddrs(pn.ID(), pn.NodeInfo.Addrs, peerstore.PermanentAddrTTL)
}

func (r *Router) DelNode(pn *tool.PeerNode) {
	r.host.Peerstore().ClearAddrs(pn.ID())
	r.host.Peerstore().RemovePeer(pn.ID())
}

func (r *Router) AllNodes() []*tool.PeerNode {

	nodes := make([]*tool.PeerNode, 0, 10)

	for _, v := range r.host.Peerstore().Peers() {
		info := r.host.Peerstore().PeerInfo(v)
		addr, _ := peer.AddrInfoToP2pAddrs(&info)
		node := tool.NewPeerNode(&info, addr[0])
		nodes = append(nodes, node)
	}
	return nodes
}

func (r *Router) Sum() int {
	return r.host.Peerstore().Peers().Len()
}

func (r *Router) Update(table Table) {

}

func (r *Router) RawData() []byte {
	data := bytes.Buffer{}
	for i := 1; i <= 256; i++ {
		//if v, _ := r.table.Load(i); v != nil {
		//	l := v.(*list.List)
		//	data.WriteString(fmt.Sprintf("%d:", i))
		//	for e := l.Front(); e != nil; e = e.Next() {
		//		addrStr := e.Value.(*tool.PeerNode).String()
		//		data.WriteString(addrStr + ";")
		//	}
		//	data.WriteString("||")
		//}
	}
	return data.Bytes()
}

// ParseData Transform the raw data into designated type, Table
func (r *Router) ParseData(raw string) Table {
	// The distances of addresses in every row are the same
	table := make(Table)
	distList := strings.Split(raw, "||")
	for _, str := range distList {
		str = strings.Trim(str, " ")
		if str == "" || str == "\n" {
			continue
		}
		row := strings.Split(str, ":")
		addrs := strings.Split(row[1], ";")
		dist, _ := strconv.ParseInt(row[0], 10, 64)
		addrList := list.New()
		for _, addr := range addrs {
			if addr == "" {
				continue
			}
			addrList.PushBack(tool.ParsePeerNode(addr))
		}
		table[int(dist)] = addrList
	}
	return table
}
