package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	for i := 1; ; i++ {
		go show(i)
		gorouteCount := runtime.NumGoroutine()
		fmt.Println(gorouteCount)
		time.Sleep(time.Second)
	}

	select {}
}

func show(index interface{}) {
	for {
		//log.Println("=================== ", index)
		time.Sleep(time.Second)
	}
}
