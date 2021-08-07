package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"time"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

// implements interface
func (dn *discoveryNotifee) HandlePeerFound(peerAddrInfo peer.AddrInfo) {
	dn.PeerChan <- peerAddrInfo
}

// initialize the mdns service
func initMDNS(ctx context.Context, peerHost host.Host, rendezvous string) chan peer.AddrInfo {
	mdnsService, err := discovery.NewMdnsService(ctx, peerHost, time.Second, rendezvous)
	if err != nil {
		panic(err)
	}
	// register with the service so that we get notifeed about peer discovery
	dn := &discoveryNotifee{
		PeerChan: make(chan peer.AddrInfo),
	}
	mdnsService.RegisterNotifee(dn)
	return dn.PeerChan
}
