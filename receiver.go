package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func syslog_receiver(out chan Message) {
	defer close(out)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	server, err := net.ListenUDP("udp", &net.UDPAddr{Port: 4242})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	dat := make([]byte, 255)
	prev_dat := make([]byte, 255)

	msg := Message{}

receive_loop:
	for {
		select {
		case _ = <-sigs:
			log.Print("Received interrupt.  Exiting...")
			break receive_loop
		default:
		}

		server.SetDeadline(time.Now().Add(time.Second))
		n, _, err := server.ReadFrom(dat)
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		if bytes.Equal(dat, prev_dat) {
			continue
		}
		copy(prev_dat, dat)
		i := 0
		for ; i < n; i++ {
			if dat[i] == '{' {
				break
			}
		}

		if i == n {
			log.Print("invalid message received")
			continue
		}

		err = json.Unmarshal(dat[i:n], &msg)
		if err != nil {
			log.Print(err)
			continue
		}

		msg.Temperature = (9.0 * msg.Temperature / 5.0) + 32.0

		out <- msg
	}
}