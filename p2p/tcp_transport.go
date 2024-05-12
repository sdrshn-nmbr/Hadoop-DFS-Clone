package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	// this mutex will PROTECT the peers map below it (common practice in go)
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func newTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

