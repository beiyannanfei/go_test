package main

import (
	"time"
	"fmt"
)

//使用select的时候，有时需要超时处理, 其中的timeout信道相当有趣

func main() {
	c1 := make(chan int)

	go func() {
		c1 <- 10
	}()

	timeout := time.After(time.Second) //timeout 是一个计时信道, 如果达到时间了，就会发一个信号出来

	for is_timeout := false; !is_timeout; {
		select {
		case v, ok := <-c1:
			fmt.Printf("c1 v: %v, ok: %v\n", v, ok)

		case v, ok := <-timeout:
			fmt.Printf("timeout v: %v, ok: %v\n", v, ok)
			is_timeout = true
		}
	}

	fmt.Println("--- finish ----")
}
