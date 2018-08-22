package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("------------------ NewTicker ------------------")
	//NewTicker 返回一个新的 Ticker，该 Ticker 包含一个通道字段，并会每隔时间段 d 就向该通道发送当时的时间。
	// 它会调整时间间隔或者丢弃 tick 信息以适应反应慢的接收者。如果d <= 0会触发panic。关闭该 Ticker 可以释放相关资源
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))

	ticker := time.NewTicker(time.Second)

	go func() {
		for tick := range ticker.C {
			fmt.Println("tick at:", tick.Format("2006-01-02 15:04:05"))
		}
	}()

	time.Sleep(time.Second * 5)
	ticker.Stop()
	fmt.Println("stopped:", time.Now().Format("2006-01-02 15:04:05"))
}
