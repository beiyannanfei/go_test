package main

import (
	"runtime"
	"fmt"
)

func main() {
	runtime.GOMAXPROCS(1)
	int_chan := make(chan int, 1)
	string_chan := make(chan string, 1)
	int_chan <- 1
	string_chan <- "hello"
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		fmt.Println(value)
	}
}

//todo select 中只要有一个case能return，则立刻执行。 当如果同一时间有多个case均能return则伪随机方式抽取任意一个执行。 如果没有一个case能return则可以执行”default”块。