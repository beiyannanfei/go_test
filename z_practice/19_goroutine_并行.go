package main

import (
	"fmt"
	"runtime"
)

//为了达到真正的并行，我们需要告诉Go我们允许同时最多使用多个核。
//回到起初的例子，我们设置最大开2个原生线程, 我们需要用到runtime包(runtime包是goroutine的调度器):

var quit = make(chan int)

func loop_1(name string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%v-%v ", name, i)
	}
	fmt.Println()
	quit <- 0
}

func loop_2(name string) {
	for i := 0; i < 10; i++ {
		runtime.Gosched() // 显式地让出CPU时间给其他goroutine
		fmt.Printf("%v-%v ", name, i)
	}
	quit <- 1
}

func main() {
	fmt.Println("======= cpu核心数量：", runtime.NumCPU())
	runtime.GOMAXPROCS(2) //最多使用2个核

	go loop_1("A") //两个goroutine会抢占式地输出数据
	go loop_1("B")

	<-quit
	<-quit

	fmt.Println("---------------------- 显式地让出CPU时间 ----------------------")

	go loop_2("AA")
	go loop_2("BB")

	<-quit
	<-quit
}
