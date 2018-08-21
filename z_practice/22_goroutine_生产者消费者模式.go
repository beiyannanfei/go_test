package main

import (
	"time"
	"fmt"
)

func main() {
	bread := make(chan int, 3)
	for i := 1; i <= 2; i++ {
		go produce(bread)
	}
	for i := 1; i <= 5; i++ {
		go consume(bread)
	}
	time.Sleep(1e9)
}

func produce(ch chan<- int) {
	for {
		ch <- 1
		fmt.Printf("produce bread len(ch): %v\n", len(ch))
		time.Sleep(100 * time.Millisecond)
	}
}

func consume(ch <-chan int) {
	for {
		<-ch
		fmt.Printf("take bread len(ch): %v\n", len(ch))
		time.Sleep(200 * time.Millisecond)
	}
}
