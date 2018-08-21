package main

/*
func main() {
	ch := make(chan int)
	<-ch //fatal error: all goroutines are asleep - deadlock! 死锁！
}
*/

/*
func main() {
	ch := make(chan int)
	ch <- 1                                  //死锁 1流入信道，堵塞当前线, 没人取走数据信道不会打开
	fmt.Println("This line code never run!") //在此行执行之前Go就会报死锁
}
*/

/*
var ch1 chan int = make(chan int)
var ch2 chan int = make(chan int)

func say(s string) {
	fmt.Println(s)
	ch1 <- <-ch2 // ch1 等待 ch2流出的数据
}

//主线等ch1中的数据流出，ch1等ch2的数据流出，但是ch2等待数据流入，两个goroutine都在等，也就是死锁
func main() {
	go say("hi")
	<-ch1 // 堵塞主线
}
*/


//总结: 非缓冲信道上如果发生了流入无流出，或者流出无流入，也就导致了死锁。或者这样理解 Go启动的所有goroutine里的非缓冲信道一定要一个线里存数据，一个线里取数据，要成对才行


