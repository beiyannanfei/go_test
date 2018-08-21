package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}

func main() {
	runtime.GOMAXPROCS(1) //Go使用单核
	go say("Hi")
	for i := 0; i < 10000000000000000; i++ { //for死循环占据了单核CPU所有的资源，而main线和say两个goroutine都在一个线程里面， 所以say没有机会执行
	}
	//解决方案
	//1. 允许Go使用多核(runtime.GOMAXPROCS)
	//2. 手动显式调动(runtime.Gosched)
}
