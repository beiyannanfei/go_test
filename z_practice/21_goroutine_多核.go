package main

import (
	"fmt"
	"runtime"
)

var quit_21 = make(chan int)

func loop_21(n string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%v ", n)
	}
	quit_21 <- 0
}

func main() {
	runtime.GOMAXPROCS(2) // 最多同时使用2个核

	nArr := [...]string{"AA", "BB", "CC"}
	for _, v := range nArr {
		go loop_21(v)
	}

	for i := 0; i < len(nArr); i++ {
		<-quit_21
	}
}

/*
➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC BB BB BB BB BB BB BB BB BB BB AA AA AA AA AA AA AA AA AA AA %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC BB BB BB BB BB BB BB BB BB BB AA AA AA AA AA AA AA AA AA AA %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC BB AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB %                                                                                                                                 ➜  z_practice git:(master) ✗ go run 21_goroutine_多核.go
CC CC CC CC CC CC CC CC CC CC AA AA AA AA AA AA AA AA AA AA BB BB BB BB BB BB BB BB BB BB %

多跑几次会看到类似这些输出(不同机器环境不一样):

执行它我们会发现以下现象:

有时会发生抢占式输出(说明Go开了不止一个原生线程，达到了真正的并行)
有时会顺序输出, 打印完AA再打印BB, 再打印CC(说明Go开一个原生线程，单线程上的goroutine不阻塞不松开CPU)
那么，我们还会观察到一个现象，无论是抢占地输出还是顺序的输出，都会有那么两个数字表现出这样的现象:

一个字符串的所有输出都会在另一个字符串的所有输出之前
原因是， 3个goroutine分配到至多2个线程上，就会至少两个goroutine分配到同一个线程里，单线程里的goroutine 不阻塞不放开CPU, 也就发生了顺序输出
*/