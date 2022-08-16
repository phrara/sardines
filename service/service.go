package service

import (
	"context"
	"sardines/router"
	"sardines/tool"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

var serv *Service

const (
	CHAT = "/chat"
	JOIN = "/join"
	RD   = "/rd"
	FT   = "/ft"
)

func init() {
	serv = &Service{
		Host:        nil,
		router:      nil,
		pingService: nil,
	}
}

type Service struct {
	Host        host.Host
	router      *router.Router
	pingService *ping.PingService
}

func NewService(host host.Host, r *router.Router) *Service {
	serv.Host = host
	serv.router = r
	serv.pingService = ping.NewPingService(host)
	return serv
}

func (s *Service) ServiceHandlerRegister() *Service {
	s.Host.SetStreamHandler(CHAT, ChatHandler)
	s.Host.SetStreamHandler(ping.ID, s.pingService.PingHandler)
	s.Host.SetStreamHandler(JOIN, JoinApplyHandler)
	s.Host.SetStreamHandler(RD, RouterDistributeHandler)
	s.Host.SetStreamHandler(FT, RecvFileHandler)
	return s
}

func (s *Service) Ping(pn *tool.PeerNode) <-chan ping.Result {
	s.Host.Peerstore().AddAddrs(pn.ID(), pn.NodeInfo.Addrs, peerstore.PermanentAddrTTL)
	return s.pingService.Ping(context.Background(), pn.ID())
}
