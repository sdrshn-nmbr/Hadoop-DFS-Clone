package main

import (
	// "fmt"
	"fmt"
	"log"

	"gohdfsclone/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshake,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpOpts)
	
	go func() {
		for {
			msg := <- tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}

}
