package main

import (
	"bytes"
	"log"
	"net"
)

func main() {
	server, err := net.ListenUDP("udp", &net.UDPAddr{Port: 4242})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	msg := make([]byte, 255)
	prev_msg := make([]byte, 255)
	for {
		n, _, err := server.ReadFrom(msg)
		if err != nil {
			log.Fatal(err)
		}
		if bytes.Equal(msg, prev_msg) {
			continue
		}
		copy(prev_msg, msg)
		i := 0
		for ; i < n; i++ {
			if msg[i] == '{' {
				break
			}
		}

		if i == n {
			log.Print(string(msg[:n]))
			log.Print("invalid message received")
			continue
		}

		log.Print(string(msg[i-1 : n]))
	}
}
