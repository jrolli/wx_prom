package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/NYTimes/gziphandler"
)

const (
	data_current = iota
)

type DataRequest struct {
	RequestType  int
	ResponseChan chan map[int]Message
}

func toFahrenheit(temp float64) float64 {
	return (9.0 * temp / 5.0) + 32.0
}

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

	requests := make(chan DataRequest)

	go http_handler(requests)

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
			switch data_req.RequestType {
			case data_current:
				data_req.ResponseChan <- data
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
			toFahrenheit(msg.Temperature),
			msg.Humidity,
			msg.Battery)

	}
	log.Print("Server exiting...")
}

func get_text_handler(requests chan DataRequest) http.Handler {
	responses := make(chan map[int]Message)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- DataRequest{data_current, responses}
		data := <-responses

		resp := ""

		chans := make([]int, 0)
		for ch := range data {
			chans = append(chans, ch)
		}

		sort.Ints(chans)

		for _, ch := range chans {
			resp = fmt.Sprintf("%s%15s: %6.1f F | %3d%% (%s)\n",
				resp,
				data[ch].Room,
				toFahrenheit(data[ch].Temperature),
				data[ch].Humidity,
				data[ch].Time[:16])
		}

		w.Write([]byte(resp))
	})
}

func get_json_handler(requests chan DataRequest) http.Handler {
	responses := make(chan map[int]Message)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- DataRequest{data_current, responses}
		data := <-responses

		outbound := []Message{}
		for _, msg := range data {
			outbound = append(outbound, msg)
		}

		resp, err := json.Marshal(outbound)
		if err != nil {
			log.Print("could not create json")
		}

		w.Write([]byte(resp))
	})
}

func http_handler(requests chan DataRequest) {
	gz := gziphandler.GzipHandler

	http.Handle("/", gz(http.FileServer(http.Dir("site"))))
	http.Handle("/api/temp.txt", gz(get_text_handler(requests)))
	http.Handle("/api/temp.json", gz(get_json_handler(requests)))

	listen_addr := ":8080"

	log.Printf("Listening at %s", listen_addr)
	log.Print(http.ListenAndServe(listen_addr, nil))
}
