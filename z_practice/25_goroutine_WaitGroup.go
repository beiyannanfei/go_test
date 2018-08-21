package main

import (
	"sync"
	"fmt"
	"time"
)

//等待所有任务退出主程序再退出

func calc_25(w *sync.WaitGroup, i int) {
	time.Sleep(time.Second * time.Duration(i))
	fmt.Println("calc: ", i)
	w.Done()
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go calc_25(&wg, i)
	}
	wg.Wait()
	fmt.Println("all goroutine finish!")
}
