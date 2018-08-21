package main

import (
	"fmt"
	"time"
)

func rand_41() chan int {
	ch := make(chan int)

	go func() {
		for {
			select { //select会尝试执行各个case, 如果都可以执行，那么随机选一个执行
			case ch <- 0:
			case ch <- 1:
			case ch <- 2:
			case ch <- 3:
			case ch <- 4:
			case ch <- 5:
			case ch <- 6:
			case ch <- 7:
			case ch <- 8:
			case ch <- 9:
			}
		}
	}()

	return ch
}

func main() {
	generator := rand_41() //初始化一个01随机生成器

	//测试，打印10个随机01
	for i := 0; i < 10; i++ {
		fmt.Println(<-generator)
		time.Sleep(time.Duration(1000000000))
	}
}
