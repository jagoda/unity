package testutils_test

import (
	. "github.com/jagoda/unity/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	
	"golang.org/x/crypto/ssh"
)

var _ = Describe("SshServer", func() {
	const (
		PASSWORD = "tester"
		USERNAME = "testy"
	)
	
	IsAServer(func () BasicServer {
		return NewSshServer(USERNAME, PASSWORD)
	})
	
	Describe("Authenticating a client", func() {
		var (
			server *SshServer
		)
		
		BeforeEach(func() {
			server = NewSshServer(USERNAME, PASSWORD)
		})
		
		AfterEach(func() {
			server.Close()
		})
		
		Context("with correct credentials", func() {
			var (
				client *ssh.Client
				err error
			)
			
			BeforeEach(func() {
				config := NewClientConfig(USERNAME, PASSWORD)
				client, err = ssh.Dial("tcp", server.Addr(), config)
			})
			
			It("allows the connection", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(client).NotTo(BeNil())
			})
		})
		
		Context("with incorrect credentials", func() {
			var (
				err error
			)
			
			BeforeEach(func() {
				config := NewClientConfig(USERNAME, PASSWORD + "z")
				_, err = ssh.Dial("tcp", server.Addr(), config)
			})
			
			It("rejects the connection", func() {
				Expect(err).To(MatchError(ContainSubstring("unable to authenticate")))
			})
		})
	})
	
	Describe("Opening a session channel", func() {
		var (
			client *ssh.Client
			err error
			server *SshServer
			session *ssh.Session
		)
		
		BeforeEach(func() {
			server = NewSshServer(USERNAME, PASSWORD)
			
			config := NewClientConfig(USERNAME, PASSWORD)
			client, err = ssh.Dial("tcp", server.Addr(), config)
			Expect(err).NotTo(HaveOccurred())
			
			session, err = client.NewSession()
		})
		
		AfterEach(func() {
			session.Close()
			server.Close()
		})
		
			
		It("creates a new session", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(session).NotTo(BeNil())
		})
	})
	
	Describe("Opening a connection channel", func() {
		var (
			client *ssh.Client
			err error
			server *SshServer
		)
		
		BeforeEach(func() {
			server = NewSshServer(USERNAME, PASSWORD)
			
			config := NewClientConfig(USERNAME, PASSWORD)
			client, err = ssh.Dial("tcp", server.Addr(), config)
			Expect(err).NotTo(HaveOccurred())
			
			_, err = client.Dial("tcp", "127.0.0.1:8080")
		})
		
		It("rejects the channel", func() {
			Expect(err).To(MatchError(ContainSubstring("only session channels are allowed")))
		})
	})
})

func NewClientConfig(username, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
}
