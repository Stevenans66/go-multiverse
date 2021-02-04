package rpc

import (
	"net"
	"net/rpc"

	"github.com/multiverse-vcs/go-multiverse/peer"
)

// connect starts an rpc server and returns a connected client.
func connect(client *peer.Client) (*rpc.Client, error) {
	service := Service{
		client: client,
	}

	server := rpc.NewServer()
	if err := server.Register(&service); err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		return nil, err
	}
	go server.Accept(listener)

	return rpc.Dial("tcp", listener.Addr().String())
}
