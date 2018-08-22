package main

import (
	"fmt"
	"time"
	"sync"
)

func main() {
	fmt.Println("--------------- Tick ---------------")
	//Tick 是 NewTicker 的封装，只提供对 Ticker 的通道的访问。如果不需要关闭 Ticker，本函数就很方便
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))

	count := 0
	wait := sync.WaitGroup{}
	wait.Add(1)
	ticker := time.Tick(time.Second)

	go func() {
		for tick := range ticker {
			count += 1
			fmt.Println("tick at:", tick.Format("2006-01-02 15:04:05"))

			if count >= 5 {
				wait.Done()
				break
			}
		}
	}()

	wait.Wait()

	fmt.Println("stopped:", time.Now().Format("2006-01-02 15:04:05"))
}
