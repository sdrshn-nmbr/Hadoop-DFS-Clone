package p2p

// * Peer is an interface that represents the remote node (the participant in the network)
type Peer interface {
	Close() error
}

// * Transport is anything that handles the communication between the nodes in the network
// ! This can be of the form TCP, UDP, Websockets, etc.
type Transport interface {
	ListenAndAccept() error
	
	// Func called consume that returns a channel of RPC
	Consume() <-chan RPC
}
