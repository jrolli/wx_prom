package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
)

type Message struct {
	Time 		string	`json:"time"`
	Model 		string	`json:"model"`
	Id 			int		`json:"id"`
	Channel 	int		`json:"channel"`
	Temperature	float64	`json:"temperature_C"`
	Humidity 	int		`json:"humidity"`
}

func main() {
	server, err := net.ListenUDP("udp", &net.UDPAddr{Port: 4242})
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	dat := make([]byte, 255)
	prev_dat := make([]byte, 255)

	msg := Message{}
	for {
		n, _, err := server.ReadFrom(dat)
		if err != nil {
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

		msg.Temperature = (9.0*msg.Temperature/5.0) + 32.0

		log.Printf("Channel %d: %10.1f F | %4d%%", msg.Channel, msg.Temperature, msg.Humidity)
	}
}
