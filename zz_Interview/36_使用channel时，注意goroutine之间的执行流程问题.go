package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	go setData(ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

func setData(ch chan string) {
	ch <- "test"
	time.Sleep(time.Second)
	ch <- "hello world"
	time.Sleep(time.Second)
	ch <- "123"
	time.Sleep(time.Second)
	ch <- "456"
	time.Sleep(time.Second)
	ch <- "789"
}

//todo 一个基于无缓存channel的发送或者取值操作，会导致当前goroutine阻塞，一直等待到另外的一个goroutine做相反的取值或者发送操作以后，才会正常跑。
// 主goroutine等待接收，另外的那一个goroutine发送了“test”并等待处理；完成通信后，打印出”test”；两个goroutine各自继续跑自己的。
// 主goroutine等待接收，另外的那一个goroutine发送了“hello world”并等待处理；完成通信后，打印出”hello world”；两个goroutine各自继续跑自己的。
// 主goroutine等待接收，另外的那一个goroutine发送了“123”并等待处理；完成通信后，打印出”123”；两个goroutine各自继续跑自己的。
// 主goroutine等待接收，另外的那一个goroutine发送了“456”并等待处理；完成通信后，打印出”456”；两个goroutine各自继续跑自己的。
// 主goroutine等待接收，另外的那一个goroutine发送了“789”并等待处理；完成通信后，打印出”789”；两个goroutine各自继续跑自己的。
// 记住：Golang的channel是用来goroutine之间通信的，且通信过程中会阻塞。
