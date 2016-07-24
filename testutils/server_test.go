package testutils_test

import (
	. "github.com/jagoda/unity/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"errors"
	"io"
	"net"
)

const (
	BIND_FAILED = "Failed to listen on port"
)

var _ = Describe("Server", func() {
	IsAServer(func () BasicServer {
		return NewServer(nil)
	})
	
	Describe("Receiving connections", func() {
		Context("when no handler is present", func() {
			var (
				connection1, connection2 net.Conn
				server *Server
			)
			
			BeforeEach(func() {
				server = NewServer(nil)
				
				var err error
				connection1, err = net.Dial("tcp", server.Addr())
				Expect(err).NotTo(HaveOccurred())
				connection2, err = net.Dial("tcp", server.Addr())
				Expect(err).NotTo(HaveOccurred())
			})
			
			AfterEach(func() {
				connection1.Close()
				connection2.Close()
				server.Close()
			})
			
			It("closes the connections", func() {
				buffer := make([]byte, 5)
				read, err := connection1.Read(buffer)
				Expect(read).To(BeZero())
				Expect(err).To(Equal(io.EOF))
				
				read, err = connection2.Read(buffer)
				Expect(read).To(BeZero())
				Expect(err).To(Equal(io.EOF))
			})
		})
		
		Context("when a handler is defined", func() {
			const (
				MESSAGE = "hello"
			)
			
			var (
				connection1, connection2 net.Conn
				server *Server
			)
			
			BeforeEach(func() {
				server = NewServer(func (connection net.Conn) {
					defer GinkgoRecover()
					
					written, err := connection.Write([]byte(MESSAGE))
					Expect(written).To(Equal(len(MESSAGE)))
					Expect(err).NotTo(HaveOccurred())
				})
				
				var err error

				connection1, err = net.Dial("tcp", server.Addr())
				Expect(err).NotTo(HaveOccurred())
				connection2, err = net.Dial("tcp", server.Addr())
				Expect(err).NotTo(HaveOccurred())
			})
			
			AfterEach(func() {
				connection1.Close()
				connection2.Close()
				server.Close()
			})
			
			It("uses the handler to service the requests", func() {
				size := 5
				buffer := make([]byte, size)
				read, err := connection1.Read(buffer)
				Expect(read).To(Equal(size))
				Expect(buffer).To(Equal([]byte(MESSAGE)))
				Expect(err).NotTo(HaveOccurred())
				
				read, err = connection2.Read(buffer)
				Expect(read).To(Equal(size))
				Expect(buffer).To(Equal([]byte(MESSAGE)))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})

type BasicServer interface {
	Addr() string
	Close()
}

func IsAServer(createServer func() BasicServer) {
	Describe("Starting a new server", func() {
		Context("when a port can be bound", func() {
			var (
				server BasicServer
			)
			
			BeforeEach(func() {
				server = createServer()
			})
			
			AfterEach(func() {
				server.Close()
			})
			
			It("creates a new server object", func() {
				Expect(server).NotTo(BeNil())
			})
			
			It("binds to a local port", func() {
				address := server.Addr()
				Expect(address).NotTo(BeEmpty())
				
				connection, err := net.Dial("tcp", address)
				Expect(err).NotTo(HaveOccurred())
				Expect(connection.Close()).To(Succeed())
			})
		})
		
		Context("when a port cannot be bound", func() {
			var (
				network *testNetwork
			)
			
			BeforeEach(func() {
				network = &testNetwork{
					bindError: true,
					network: DefaultNetwork,
				}
				DefaultNetwork = network
			})
			
			AfterEach(func() {
				DefaultNetwork = network.network
			})
			
			It("panics", func() {
				start := func () {
					createServer()
				}
				
				Expect(start).To(Panic())
			})
		})
	})
	
	Describe("Stopping a server", func() {
		var (
			server BasicServer
		)
		
		BeforeEach(func() {
			server = createServer()
			server.Close()
		})
		
		It("causes it to stop accepting connections", func() {
			address := server.Addr()
			Expect(address).NotTo(BeEmpty())
			
			_, err := net.Dial("tcp", address)
			Expect(err).To(MatchError(ContainSubstring("refused")))
		})
	})
}

type testNetwork struct {
	bindError bool
	network Network
}

func (network *testNetwork) Listen() (net.Listener, error) {
	if network.bindError {
		return nil, errors.New(BIND_FAILED)
	}
	
	return network.network.Listen()
}