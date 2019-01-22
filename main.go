package main

import (
	"log"
)

type Message struct {
	Time        string  `json:"time"`
	Model       string  `json:"model"`
	Id          int     `json:"id"`
	Channel     int     `json:"channel"`
	Temperature float64 `json:"temperature_C"`
	Humidity    int     `json:"humidity"`
	Battery     string  `json:"battery"`
}

func main() {
	ch := make(chan Message)

	log.Print("Starting UDP receiver...")
	go syslog_receiver(ch)

	log.Print("Starting HTTP server...")

	log.Print("Started")
}
