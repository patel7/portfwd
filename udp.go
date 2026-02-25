package main

import (
	"log"
	"net"

	"github.com/libp2p/go-reuseport"
)

func udpForward(forward ForwardStruct) {
	src, err := reuseport.ListenPacket(forward.Protocol, forward.From)
	if err != nil {
		log.Fatal("Failed to listen on UDP port/address %s : %v", forward.From, err)
	}
	defer src.Close()

	dstAddr, err := net.ResolveUDPAddr(forward.Protocol, forward.To)
	if err != nil {
		log.Fatal("Error resolving destination address: %v\n", err)
	}

	dst, err := net.DialUDP(forward.Protocol, nil, dstAddr)
	if err != nil {
		log.Fatal("Outgoing UDP connection failed: %v", err)
	}
	dst.Close()

	for {
		buf := make([]byte, 65535)
		n, _, err := src.ReadFrom(buf)
		if err != nil {
			log.Fatal("Error reading from UDP socket: %v\n", err)
		}

		_, _ = dst.Write(buf[:n])
	}
}
