package main

import (
	"github.com/highras/fpnn-sdk-go/src/fpnn"
	"time"
	"log"
)

func main() {
	endpoint := "127.0.0.1:12321"
	client := fpnn.NewTCPClient(endpoint)

	client.SetOnConnectedCallback(func(connId uint64) {
		log.Printf("Connected to %s, connId is %d\n", endpoint, connId)
	})

	client.SetOnClosedCallback(func(connId uint64) {
		log.Printf("Connection no %s, closed, connId is: %d\n", endpoint, connId)
	})

	if ok := client.Connect(); !ok {
		log.Printf("connect to %s failed.", endpoint)
	}

	defer client.Close()

	time.Sleep(time.Second) //-- Waiting for the closed event printed.
}
