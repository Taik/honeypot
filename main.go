package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/taik/honeypot/ssh"
)

var port = 22

func init() {
	flag.IntVar(&port, "port", 22, "the default ssh port to listen to")
	flag.Parse()
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("unable to bind to port: %s\n", err)
	}
	defer l.Close()

	log.Printf("listening on port %d", port)

	sshServer, err := ssh.NewServer(l)
	if err != nil {
		log.Fatalf("unable to create new SSH server: %s\n", err)
	}
	go sshServer.Start()

	for attempts := range sshServer.AttemptsCh() {
		log.Printf("%v\n", attempts)
	}
}
