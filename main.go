package main

import (
	"log"
	"net"

	"github.com/taik/honeypot/ssh"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatalf("unable to bind to port: %s\n", err)
	}
	defer l.Close()

	sshServer, err := ssh.NewServer(l)
	if err != nil {
		log.Fatalf("unable to create new SSH server: %s\n", err)
	}
	go sshServer.Start()

	for attempts := range sshServer.AttemptsCh() {
		log.Printf("%v\n", attempts)
	}
}
