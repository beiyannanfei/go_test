package main

import (
	"time"
	"fmt"
)

func main() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	timer1 := time.NewTimer(time.Second * 2)
	res := <-timer1.C
	fmt.Printf("Timer 1 expired. res: %v, res: %v\n", res.Unix(), res.Format("2006-01-02 15:04:05"))

	fmt.Println("=========================================================================")

	timer2 := time.NewTimer(time.Second)
	go func() {
		res2 := <-timer2.C
		fmt.Printf("Timer2 expired. res2: %v, res2: %v\n", res2.Unix(), res2.Format("2006-01-02 15:04:05"))
	}()

	stop2 := timer2.Stop()
	if stop2 {
		fmt.Printf("Timer2 stopped. stop2: %v\n", stop2)
	}

	fmt.Println("=========================================================================")

	ticker := time.NewTicker(time.Second) //ticker是一个定时触发的计时器，它会以一个间隔(interval)往Channel发送一个事件(当前时间)
	go func() {
		for t := range ticker.C {
			fmt.Println("ticker at: ", t.Format("2006-01-02 15:04:05"))
		}
	}()

	time.Sleep(time.Second * 10)
	ticker.Stop()	//ticker也可以通过Stop方法来停止。一旦它停止，接收者不再会从channel中接收数据了
}
