package p2p

import (
	// "bytes"
	"fmt"
	"net"
	"sync"
)

// ! represents the remote node over a TCP established connection
type TCPPeer struct {
	// * conn is the underlying connection of the peer node
	conn net.Conn

	// * if we dial and retrieve a connection: outbound == true
	// * else if we accept and retrieve a connection: outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// * Close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC

	// * this mutex will PROTECT the peers map below it (common practice in go)
	// * aka the map can only be accessed one goroutine at a time
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

// ! constructor method that creates the TCPTransport object
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// * Consume implements the transport interface which will return a read-only channel for reading the incoming msgs received from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// * gonna be called in the server
func (t *TCPTransport) ListenAndAccept() error {

	// * listening
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	// * accepting
	go t.StartAcceptLoop() // *-> goroutine using method of TCPTransport pointer

	// * no error when listening and accepting
	return nil
}

// ! method of TCPTransport
func (t *TCPTransport) StartAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept error: %s\n", err)
		}
		fmt.Printf("New incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()

		t.rpcch <-  rpc
	}
}
