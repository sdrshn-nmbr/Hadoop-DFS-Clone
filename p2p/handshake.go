package p2p

import (
	"net"
)

type HandShaker interface {
	HandShake() net.Error
}