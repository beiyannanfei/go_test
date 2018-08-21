package main

import (
	"fmt"
	"time"
)

func main() {
	intChan := make(chan interface{}, 10)

	for i := 0; i < 10; i++ {
		intChan <- i
	}

	close(intChan)

	for {
		i, ok := <-intChan
		fmt.Println(i, ok)
		if !ok {
			fmt.Println("channel is close.")
			break
		}
		time.Sleep(time.Second)
	}
}
