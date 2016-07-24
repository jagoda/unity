package testutils

import (
	"golang.org/x/crypto/ssh"
	
	"fmt"
	"net"
)

const (
	HOST_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCj1/fBi6O3rVyqODOCiddjLHvyYZzqtg/atoeqgZHN4ZoI6THl
F/gjwLljCtgK30YisLJ0etsFMX04nWRhTtMFpTgiGjLv+2pwtj8tM8ibBrNMuIU6
ExcTScKC1XO++zb2KAjB8/YjM7qqFa5Fe68HqYPFZ/TGowbed6cg+I0H7QIDAQAB
AoGARLkUqbEJpcRoptlV+b4Zgvaospzw9Z4R/lorg4A/iQzE0FIH5KDuKwjuebWV
stG+GFTnNWZGseH1NScHcK1gY0TxyTfLl1UmFnGl3rqFeY1vWlSNScWXO5lvrQdS
Hrd6OnIgTkDMDXProZABJgVi4gi1nXVMB2yy6bkM6HdL+qECQQDSDvZ9RmM/kf67
3bGA4pKdlNPxm8DemJerZJLl0HHxqf0QnmRRzleN99UZz7NhXz8XvaCx/kSHHwEE
aKLYMAllAkEAx613eqUz/Cysk6XwLgyiP/dxp1EyU4r/Xr+f1ip8RC9HCVPYBYBx
HizMSgGyy4t0wgSE4pN6dNcIaRi5cnhf6QJBAIldYZFwuyxyI5vlENPQ1sIX9jWU
fh7SuLRLM8j3c9vuJuB8Q+w/PIBJZmDXs11ktNSom/Jp7ZRTEQ46AuvPFgUCQHbi
xTx0mBsQzj+qkPIZ+5ByV2zzXy92ls1m8lelPA+sxnK7ROchrjf1HD0D/dxKz92k
qynr/QEL9qn8Wo3ZNeECQECwi/EeUWIT90+CKjY8vilFDC952xJhCEkB9vbQz2YW
vPJnR8euf24Ec6XfoUDei2Mz3XGI1qTYc0vsJO+IWIA=
-----END RSA PRIVATE KEY-----`
)

// SshServer provides an SSH server implementation that is suitable for testing clients.
type SshServer struct {
	*Server
	
	config *ssh.ServerConfig
}

// NewSshServer creates a new SshServer instance.
func NewSshServer(username, password string) *SshServer {
	sshServer := new(SshServer)
	
	sshServer.config = &ssh.ServerConfig {
		PasswordCallback: func(context ssh.ConnMetadata, clientPassword []byte) (*ssh.Permissions, error) {
			if context.User() == username && string(clientPassword) == password {
				return nil, nil
			}
			
			return nil, fmt.Errorf("Failed to authenticate user '%s'", context.User())
		},
	}
	
	hostKey, _ := ssh.ParsePrivateKey([]byte(HOST_KEY))
	sshServer.config.AddHostKey(hostKey)
	
	sshServer.Server = NewServer(sshServer.handshake)
	return sshServer
}

func (server *SshServer) handshake(connection net.Conn) {
	serverConnection, channels, requests, err := ssh.NewServerConn(connection, server.config)
	
	if err != nil {
		return
	}
	
	defer serverConnection.Close()
	go ssh.DiscardRequests(requests)
	
	for newChannel := range channels {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "only session channels are allowed")
			continue
		}
		
		channel, requests, _ := newChannel.Accept()
		go server.session(channel, requests)
	}
}

func (server *SshServer) session(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()
	
	for request := range requests {
		request.Reply(true, nil)
	}
}