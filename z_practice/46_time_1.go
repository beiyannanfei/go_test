package main

import (
	"time"
	"fmt"
	"sync"
)

func main() {
	t, err := time.Parse("2006-01-02 15:04:05", "2018-08-21 19:30:50")
	fmt.Println("time.parse:", t, err) //time.parse: 2018-08-21 19:30:50 +0000 UTC <nil>

	fmt.Println("------------------ NewTimer ------------------")
	//NewTimer 创建一个 Timer，它会在最少过去时间段 d 后到期，向其自身的 C 字段发送当时的时间
	fmt.Println("now:", time.Now().Format("2006-01-02 15:04:05"))
	timer := time.NewTimer(time.Second * 2)
	afterTime, ok := <-timer.C
	fmt.Printf("afterTime: %v, ok: %v\n", afterTime.Format("2006-01-02 15:04:05"), ok)

	fmt.Println("----------------------- AfterFunc -----------------------")
	//AfterFunc 另起一个 go 协程等待时间段 d 过去，然后调用 f。它返回一个 Timer，可以通过调用其 Stop 方法来取消等待和对 f 的调用。
	wait := sync.WaitGroup{}
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))
	wait.Add(1)
	timer = time.AfterFunc(time.Second*3, func() {
		fmt.Println("get timer:", time.Now().Format("2006-01-02 15:04:05"))
		wait.Done()
	})

	fmt.Println("finish:", time.Now().Format("2006-01-02 15:04:05"))
	timer.Reset(2 * time.Second) //Reset 使 t 重新开始计时，（本方法返回后再）等待时间段 d 过去后到期。如果调用时 t 还在等待中会返回真；如果 t 已经到期或者被停止了会返回假
	//timer.Stop()
	wait.Wait()

	fmt.Println("--------------------- Reset ---------------------")
	//Reset 使 t 重新开始计时，（本方法返回后再）等待时间段 d 过去后到期。如果调用时 t 还在等待中会返回真；如果 t 已经到期或者被停止了会返回假。
	wait = sync.WaitGroup{}
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))

	wait.Add(1)
	timer = time.NewTimer(time.Second * 2)

	go func() {
		<-timer.C

		fmt.Println("get timer:", time.Now().Format("2006-01-02 15:04:05"))
		wait.Done()
	}()

	time.Sleep(time.Second)
	fmt.Println("sleep:", time.Now().Format("2006-01-02 15:04:05"))

	timer.Reset(time.Second * 3) //从现在开始重新计时
	wait.Wait()
	fmt.Println("finish:", time.Now().Format("2006-01-02 15:04:05"))

	fmt.Println("--------------------- Stop ---------------------")
	//Stop 停止 Timer 的执行。如果停止了 t 会返回真；如果 t 已经被停止或者过期了会返回假。Stop 不会关闭通道 t.C，以避免从该通道的读取不正确的成功。
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))
	timer = time.NewTimer(time.Second * 2)

	go func() {
		<-timer.C
		fmt.Println("get timer")
	}()

	time.Sleep(time.Second)
	if timer.Stop() {
		fmt.Println("timer stopped:", time.Now().Format("2006-01-02 15:04:05"))
	}
}
