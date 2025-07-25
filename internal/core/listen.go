package core

import (
	"net"
)

func Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
