package main

import "fmt"

func main() {
	ch := make(chan int)
	ch <- 1 //fatal error: all goroutines are asleep - deadlock!
	//todo chan一旦被写入数据后，当前goruntine就会被阻塞，知道有人接收才可以（即 “ <- ch”），如果没人接收，它就会一直阻塞着。而如果chan带一个缓冲，就会把数据放到缓冲区中，直到缓冲区满了，才会阻塞
	fmt.Println("success")
}
