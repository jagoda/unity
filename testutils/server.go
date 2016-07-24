package testutils

import (
	"net"
)

var (
	// DefaultNetwork allows the low level network API to be abstracted for test purposes.
	DefaultNetwork Network
)

func init () {
	DefaultNetwork = new(defaultNetwork)
}

// Handler is a function that handles incoming connections.
type Handler func (net.Conn)

// Network is an API abstraction for network interactions.
type Network interface {
	Listen() (net.Listener, error)
}

// Server provides utilities for managing incoming network connections.
type Server struct {
	done chan struct{}
	listener net.Listener
}

// NewServer creates a new Server instance. The new instance automatically starts listening for connections.
func NewServer(handler Handler) *Server {
	listener, err := DefaultNetwork.Listen()
	
	if err != nil {
		panic(err)
	}
	
	server := &Server {
		done: make(chan struct{}),
		listener: listener,
	}
	
	go server.waitForConnections(handler)
	return server
}

// Addr returns the local address that the server is listening on.
func (server *Server) Addr() string {
	return server.listener.Addr().String()
}

// Close causes the server to stop accepting new connections.
func (server *Server) Close() {
	server.listener.Close()
}

func (server *Server) waitForConnections(handler Handler) {
	var (
		connection net.Conn
		err error
	)
	
	for {
		connection, err = server.listener.Accept()
		
		if err != nil {
			break
		}
		
		if handler != nil {
			go handler(connection)
		} else {
			connection.Close()
		}
	}
}

type defaultNetwork struct {}

func (network *defaultNetwork) Listen() (net.Listener, error) {
	address := &net.TCPAddr{}
	return net.ListenTCP("tcp", address)
}