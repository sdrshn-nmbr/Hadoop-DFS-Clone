package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(test *testing.T) {
	listenAddr := ":4000"
	t := NewTCPTransport(listenAddr)

	assert.Equal(test, t.listenAddress, listenAddr)

	assert.Nil(test, t.ListenAndAccept())
}
