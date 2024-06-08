package p2p

import (
	"fmt"
	"net"
	"sync"
)

// represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peer node
	conn net.Conn

	// if we dial and receive a connection -> outbount == true
	// else if we accept and receive a connection -> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	// this mutex will PROTECT the peers map below it (common practice in go) -> aka the map can only be accessed one goroutine at a time
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport { // constructor method that creates the TCPTransport object
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error { // gonna be called in the server

	// listening
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	// accepting
	go t.StartAcceptLoop() // -> goroutine using method of TCPTransport pointer

	return nil // no error when listening and accepting
}

func (t *TCPTransport) StartAcceptLoop() { // method of TCPTransport pointer
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept error: %s\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	fmt.Printf("New incoming connection %+v\n", peer)
}
