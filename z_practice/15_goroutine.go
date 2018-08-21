package main

import (
	"fmt"
	"time"
)

func loop(name string) {
	for i := 0; i < 1000; i++ {
		fmt.Printf("%v-%v  ", name, i)
	}
	fmt.Println()
}

var complete chan int = make(chan int)

func loop_chan(name string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("loop_chan-%v-%v  ", name, i)
	}
	fmt.Println()
	complete <- 0 //信道通知主线程结束
}

//数据流入无缓冲信道, 如果没有其他goroutine来拿走这个数据，那么当前线阻塞

func loop_chan1(name string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("loop_chan1-%v-%v  ", name, i)
	}
	fmt.Println()
	complete <- 0 //信道通知主线程结束	 如果complete被存入了数据并行没有被取出则会阻塞
}

func main() {
	fmt.Println("------------------ 串行执行 ------------------")
	loop("A")
	loop("B")

	fmt.Println("------------------ go执行 ------------------")
	go loop("AA")
	loop("BB")
	time.Sleep(time.Second)

	fmt.Println("------------------ 并行执行 ------------------")
	go loop("AAA")
	go loop("BBB")
	go loop("CCC")
	time.Sleep(time.Second) //如果此处不sleep的话主线就会过早跑完，loop线都没有机会执行

	fmt.Println("------------------ 信道通知 ------------------")
	go loop_chan("Z")
	go loop_chan1("Z")
	<-complete //直到线程跑完, 取到消息. main在此阻塞住
	<-complete //直到线程跑完, 取到消息. main在此阻塞住
}
