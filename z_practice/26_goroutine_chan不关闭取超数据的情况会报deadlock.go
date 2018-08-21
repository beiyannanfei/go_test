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

	for {
		//十次后 fatal error: all goroutines are asleep - deadlock!
		i := <-intChan
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
