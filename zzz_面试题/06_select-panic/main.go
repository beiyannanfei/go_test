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
	//有可能会发生异常，如果没有selct这段代码，就会出现线程阻塞，当有selct这个语句后，系统会随机抽取一个case进行判断，只有有其中一条语句正常return，此程序将立即执行。
	select {
	case value := <-int_chan:
		fmt.Println(value)
	case value := <-string_chan:
		panic(value)
	}
}
