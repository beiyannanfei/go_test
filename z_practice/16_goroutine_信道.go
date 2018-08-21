package main

import (
	"fmt"
	"time"
)

//https://blog.csdn.net/skh2015java/article/details/60330785
func main() {
	message := make(chan string) //创建信道  <=> var message chan string = make(chan string)

	go func(msg string) {
		time.Sleep(time.Second)
		fmt.Println("开始存数据...")
		message <- msg //将消息存入信道
	}("hello goroutine!")

	fmt.Println("开始取消息...")
	var info string
	info = <-message //从信道取消息  如果信道中没有数据则会一直阻塞，直到存入数据 无缓冲的信道在取消息和存消息的时候都会挂起当前的goroutine
	fmt.Println(info)

	for i := 0; i < 100; i++ { //如果上边从读取数据阻塞了这里就不会执行
		fmt.Printf("i: %v ", i)
	}
	fmt.Println()

	fmt.Println("---------------------- 无缓冲信道的数据进出顺序 ---------------------")

	//开启5个routine
	for i := 0; i < 5; i++ {
		go foo(i)
	}

	//取出信道中得数据
	for i := 0; i < 5; i++ {
		fmt.Println(<-ch)	//当数字比较小的时候输出还算有规律 4 0 1 2 3，当数字较大时无序输出
	}
}

var ch = make(chan int)

func foo(id int) {
	ch <- id
}
