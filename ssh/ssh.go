package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type UsernamePasswordAttempt struct {
	Username   string
	Password   string
	RemoteAddr net.Addr
}

func handlePassword(attemptCh chan<- *UsernamePasswordAttempt) func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	return func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		attemptCh <- &UsernamePasswordAttempt{
			Username:   conn.User(),
			Password:   string(password),
			RemoteAddr: conn.RemoteAddr(),
		}

		return nil, errors.New("invalid password")
	}
}

type Server struct {
	listener  net.Listener
	config    *ssh.ServerConfig
	attemptCh chan *UsernamePasswordAttempt
}

// NewServer returns an ssh server bootstrapped.
func NewServer(l net.Listener) (*Server, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate new RSA key")
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		errors.Wrap(err, "unable to convert RSA key to signer")
	}

	attemptCh := make(chan *UsernamePasswordAttempt)

	config := &ssh.ServerConfig{
		PasswordCallback: handlePassword(attemptCh),
		MaxAuthTries:     100,
	}
	config.AddHostKey(signer)

	return &Server{
		listener:  l,
		config:    config,
		attemptCh: attemptCh,
	}, nil
}

// Start runs the ssh server.
func (s *Server) Start() error {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			return errors.Wrap(err, "ssh: unable to accept incoming connection")
		}

		go func(conn net.Conn) {
			defer conn.Close()
			ssh.NewServerConn(c, s.config)
		}(c)
	}
}

// AttemptsCh returns the channel on which username/password auth attempts are written to.
func (s *Server) AttemptsCh() <-chan *UsernamePasswordAttempt {
	return s.attemptCh
}
