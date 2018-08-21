package main

import (
	"fmt"
	"time"
)

//多个goroutine处理任务；
//等待一组channel的返回结果。

func calc_31(taskChan chan int, resChan chan int, exitChan chan bool, taskFlag int) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("error...", err)
			return
		}
	}()

	for v := range taskChan { //任务处理逻辑
		fmt.Printf("========= calc_31 ========= v: %v, taskFlag: %v\n", v, taskFlag)
		flag := true
		for i := 2; i < v; i++ {
			if v%i == 0 {
				flag = false
				break
			}
		}

		if flag { //结果存入chan
			resChan <- v
		}
	}

	//处理完毕，存入结果到退出信道
	exitChan <- true
}

func main() {
	intChan := make(chan int, 1000) //任务信道
	resChan := make(chan int, 1000) //结果信道
	exitChan := make(chan bool, 8)  //退出信道

	go func() {
		for i := 0; i < 10; i++ {
			intChan <- i
			time.Sleep(time.Second)
		}
		close(intChan)
	}()

	//启动8个goroutine做任务
	for i := 0; i < cap(exitChan); i++ {
		go calc_31(intChan, resChan, exitChan, i)
	}

	go func() {
		//等待所有goroutine结束
		for i := 0; i < 8; i++ {
			<-exitChan
		}

		close(resChan)
		close(exitChan)
	}()

	for v := range resChan {
		fmt.Println(v)
	}
}
