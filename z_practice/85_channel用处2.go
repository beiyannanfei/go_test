package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i = i + 1 {
			c <- i
			time.Sleep(time.Second * 1)
		}
		close(c)	//如果没有这句代码，则下边的for循环会报错 fatal error: all goroutines are asleep - deadlock!
	}()

	for i := range c {
		fmt.Println(i)
	}

	fmt.Println("finished.")
}
