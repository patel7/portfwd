package main

import (
	"io"
	"log"
	"net"

	"github.com/libp2p/go-reuseport"
)

func tcpForward(forward ForwardStruct) {
	listener, err := reuseport.Listen(forward.Protocol, forward.From)

	if err != nil {
		log.Fatal("Failed to listen on TCP port/address %s : %v", forward.From, err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal("The incoming TCP connection could not be accepted: %v", err)
		}

		client, err := net.Dial(forward.Protocol, forward.To)

		if err != nil {
			log.Fatal("Outgoing TCP connection failed: %v", err)
			conn.Close()
			continue
		}

		go func() {
			defer client.Close()
			defer conn.Close()
			io.Copy(client, conn)
		}()

		go func() {
			defer client.Close()
			defer conn.Close()
			io.Copy(conn, client)
		}()
	}
}
