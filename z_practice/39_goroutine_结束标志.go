package main

import (
	"time"
	"fmt"
)

func main() {
	c, quit := make(chan int), make(chan int)

	go func() {
		for i := 0; i < 10; i++ { //每秒向信道存入一个值
			c <- i
			time.Sleep(time.Second)
		}
	}()

	go func() {
		time.Sleep(time.Second * time.Duration(5)) //五秒后向退出信道传入一个值
		quit <- 123456
	}()

	for is_quit := false; !is_quit; {
		select { // 监视信道c的数据流出
		case v, ok := <-c:
			fmt.Printf("receive value: %v, ok: %v\n", v, ok)

		case vQuit, ok := <-quit: // quit信道有输出，关闭for循环
			fmt.Printf("quit vQuit: %v, ok: %v\n", vQuit, ok)
			is_quit = true
		}
	}
}
