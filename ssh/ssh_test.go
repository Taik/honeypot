package ssh_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taik/honeypot/ssh"
	ssh2 "golang.org/x/crypto/ssh"
)

func TestServer_Start(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer l.Close()

	s, err := ssh.NewServer(l)
	require.NoError(t, err)

	go s.Start()

	go ssh2.Dial(l.Addr().Network(), l.Addr().String(), &ssh2.ClientConfig{
		User: "test-user",
		Auth: []ssh2.AuthMethod{
			ssh2.Password("test-password"),
		},
		HostKeyCallback: ssh2.InsecureIgnoreHostKey(),
	})

	select {
	case got := <-s.AttemptsCh():
		require.NotNil(t, got)
		assert.Equal(t, "test-user", got.Username)
		assert.Equal(t, "test-password", got.Password)

	case <-time.After(1 * time.Second):
		t.Fatal("timed out")
	}

}
