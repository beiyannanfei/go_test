package main

import (
	"sync"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//在go语言中优雅退出goroutines,通常需要做以下3点：
//1. 向各个goroutines发通知，令其退出，如shutdown.
//2. 等待各个goroutines都退出，如: sync.WaitGroup.
//3. 在退出goroutine之前，确保数据不丢失（1.停止生产数据。2.关闭数据channel messages. 3. 消费者goroutine检查判断数据channel messages是否有效，若无效，则退出。）

func consumer(message <-chan int, shutdown <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case msg, ok := <-message:
			if !ok {
				fmt.Println("no data, exit.")
				return
			}

			fmt.Println("msg: ", msg)

		case _ = <-shutdown:
			fmt.Println("all done!")
			return
		}
	}
}

func main() {
	shutdown := make(chan int)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println("sig: ", sig)
		shutdown <- 0 //or  close(shutdown)
	}()

	messages := make(chan int, 10)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go consumer(messages, shutdown, wg)
	for i := 0; i < 10; i++ {
		messages <- i
		time.Sleep(time.Second)
	}

	close(messages)
	fmt.Println("wait!")
	wg.Wait()
}
