package main

import (
	"time"
	"fmt"
)

func main() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			fmt.Printf("[%v] 每5s输出一次.\n", time.Now().Format("2006-01-02 15:04:05"))
		}
		ch <- 1
	}()
	<-ch
}
