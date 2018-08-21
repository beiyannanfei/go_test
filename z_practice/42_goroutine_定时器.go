package main

import (
	"time"
	"fmt"
)

//利用信道做定时器

func timer(duration time.Duration) chan bool {
	ch := make(chan bool)

	go func() {
		time.Sleep(duration)
		ch <- true //时间到了
	}()

	return ch
}

func main() {
	timeout := timer(time.Second) //定时1s

FOR:
	for {
		select {
		case <-timeout:
			fmt.Println("already 1s!")
			break FOR
		}
	}
}
