package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

const (
	data_current = iota
)

func http_server(msgs chan Message, exit chan bool) {
	rooms := map[int]string{
		1: "Outside",
		2: "Basement",
		3: "Living Room",
		4: "Kitchen",
		5: "Master Bedroom",
		6: "Spare Bedroom",
		7: "Office",
		8: "Attic",
	}

	data := make(map[int]Message)

	requests := make(chan int)
	responses := make(chan map[int]Message)

	go http_handler(requests, responses)

data_loop:
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				break data_loop
			}
			msg.Room = rooms[msg.Channel]
			data[msg.Channel] = msg
		case data_req, ok := <-requests:
			if !ok {
				log.Fatal("connection to frontend lost")
			}
			switch data_req {
			case data_current:
				responses <- data
			default:
				log.Fatal("unknown data request received")
			}
		case <-exit:
			break data_loop
		}
	}
	for ch, msg := range data {
		log.Printf("%-14s: %6.1f F | %3d%% (Batt: %s)",
			rooms[ch],
			msg.Temperature,
			msg.Humidity,
			msg.Battery)

	}
	log.Print("Server exiting...")
}

func get_handler(requests chan int, responses chan map[int]Message) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		requests <- data_current
		data := <-responses

		resp := ""

		chans := make([]int, 0)
		for ch, _ := range data {
			chans = append(chans, ch)
		}

		sort.Ints(chans)

		for _, ch := range chans {
			resp = fmt.Sprintf("%s%15s: %6.1f F | %3d%% (%s)\n",
				resp,
				data[ch].Room,
				data[ch].Temperature,
				data[ch].Humidity,
				data[ch].Time[:16])
		}

		w.Write([]byte(resp))
	})
}

func http_handler(requests chan int, responses chan map[int]Message) {
	http.HandleFunc("/temp.txt", get_handler(requests, responses))

	listen_addr := ":8080"

	log.Printf("Listening at %s", listen_addr)
	log.Print(http.ListenAndServe(listen_addr, nil))

}
