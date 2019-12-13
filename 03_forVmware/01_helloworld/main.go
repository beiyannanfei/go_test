package main

import (
	"log"
	"time"
)

func show() {
	log.Println("=============== hello world ===============")
}

func main() {
	for {
		show()
		time.Sleep(time.Second * 1)
	}
}
