package main

import (
	"log"
	"os"
	"os/signal"
	"time"
)

type Message struct {
	Time        string  `json:"time"`
	Model       string  `json:"model"`
	Id          int     `json:"id"`
	Channel     int     `json:"channel"`
	Temperature float64 `json:"temperature_C"`
	Humidity    int     `json:"humidity"`
	Battery     string  `json:"battery"`
	Room        string  `json:"room"`
	PosixTime   int64   `json:"posix_time"`
}

func main() {
	msgs := make(chan Message)
	exit := make(chan bool)

	log.Print("Starting UDP receiver...")
	go syslog_receiver(msgs, exit)

	log.Print("Starting HTTP server...")
	go http_server(msgs, exit)

	log.Print("Started")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	select {
	case <-exit:
	case <-sigs:
		log.Print("Received interrupt.  Exiting...")
		close(exit)
	}
	log.Print("Delaying main for 5 seconds...")
	time.Sleep(5 * time.Second)
	log.Print("Exited cleanly")
}
