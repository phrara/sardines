package router

import (
	"bytes"
	"container/list"
	"fmt"
	"sardines/tool"
	"strconv"
	"strings"
	"sync"
)

var rt *Router

func init() {

	rt = &Router{
		table:    sync.Map{},
		HostNode: nil,
	}

}

type Table map[int]*list.List

type Router struct {
	table    sync.Map
	HostNode *tool.PeerNode
}

func InitRouterTable(hn string) *Router {
	n := tool.ParsePeerNode(hn)
	l := list.New()
	l.PushBack(n)
	rt.table.Store(0, l)
	rt.HostNode = n
	return rt
}

func (r *Router) AddNode(pn *tool.PeerNode) {
	hnode := r.HostNode
	dist := tool.GetPeerDist(hnode.ID().String(), pn.ID().String())

	if v, _ := r.table.Load(dist); v == nil {
		l := list.New()
		l.PushBack(pn)
		r.table.Store(dist, l)
	} else {
		if b := r.ContainsAt(dist, pn); !b {
			l, _ := r.table.Load(dist)
			l.(*list.List).PushBack(pn)
		}
	}
}

func (r *Router) DelNode(pn *tool.PeerNode) {
	hnode := r.HostNode
	dist := tool.GetPeerDist(hnode.ID().String(), pn.ID().String())

	if v, _ := r.table.Load(dist); v == nil {
		return
	} else {
		l, _ := r.table.Load(dist)
		for e := l.(*list.List).Front(); e != nil; e = e.Next() {
			if e.Value.(*tool.PeerNode).String() == pn.String() {
				l.(*list.List).Remove(e)
			}
		}
	}
}

func (r *Router) Contains(pn *tool.PeerNode) bool {
	hnode := r.HostNode
	dist := tool.GetPeerDist(hnode.ID().String(), pn.ID().String())
	if dist == 0 {
		return true
	}
	if v, _ := r.table.Load(dist); v == nil {
		return false
	} else {
		l, _ := r.table.Load(dist)
		for e := l.(*list.List).Front(); e != nil; e = e.Next() {
			if e.Value.(*tool.PeerNode).String() == pn.String() {
				return true
			}
		}
	}
	return false
}

func (r *Router) ContainsAt(dist int, pn *tool.PeerNode) bool {
	if dist == 0 {
		return true
	}
	if v, _ := r.table.Load(dist); v == nil {
		return false
	} else {
		l, _ := r.table.Load(dist)
		for e := l.(*list.List).Front(); e != nil; e = e.Next() {
			if e.Value.(*tool.PeerNode).String() == pn.String() {
				return true
			}
		}
	}
	return false
}

func (r *Router) AllNodes() *list.List {
	nodes := list.New()
	for i := 1; i <= 256; i++ {
		if v, _ := r.table.Load(i); v != nil {
			l := v.(*list.List)
			for e := l.Front(); e != nil; e = e.Next() {
				nodes.PushBack(e.Value)
			}
		}
	}
	return nodes
}

func (r *Router) GetNodes(dist int) *list.List {
	nodes := list.New()
	value, ok := r.table.Load(dist)
	if !ok || value == nil {
		return nil
	}
	l := value.(*list.List)
	for i := l.Front(); i != nil; i = i.Next() {
		nodes.PushBack(i.Value)
	}
	return nodes
}

func (r *Router) Sum() int {
	size := 0
	for i := 1; i <= 256; i++ {
		l, _ := r.table.Load(i)
		if l != nil {
			size = size + l.(*list.List).Len()
		}
	}
	return size
}

func (r *Router) Clear() {
	for i := 1; i <= 256; i++ {
		r.table.Store(1, nil)
	}
}

func (r *Router) Update(table Table) {
	for i := 1; i <= 256; i++ {
		if table[i] != nil {
			for e := table[i].Front(); e != nil; e = e.Next() {
				if e.Value.(*tool.PeerNode).String() == r.HostNode.String() {
					continue
				}
				r.AddNode(e.Value.(*tool.PeerNode))
			}
		}
	}
}

/*
	Sample:
	1:/ip4/127.0.0.1/tcp/2300/p2p/Qmx;/ip4/127.0.0.1/tcp/2301/p2p/Qmx;||2:/ip4/127.0.0.1/tcp/2302/p2p/Qmx;||\n
*/
// Router table will be transformed into a type of string data, witch is raw
// This kind of string should make the consumption of router distribution lower and easier
func (r *Router) RawData() []byte {
	data := bytes.Buffer{}
	for i := 1; i <= 256; i++ {
		if v, _ := r.table.Load(i); v != nil {
			l := v.(*list.List)
			data.WriteString(fmt.Sprintf("%d:", i))
			for e := l.Front(); e != nil; e = e.Next() {
				addrStr := e.Value.(*tool.PeerNode).String()
				data.WriteString(addrStr + ";")
			}
			data.WriteString("||")
		}
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
