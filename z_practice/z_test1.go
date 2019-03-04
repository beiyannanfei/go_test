package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 1000000; i++ {
		go rou(i)
	}

	time.Sleep(time.Second * 20)
}

func rou(i int) {
	for {
		fmt.Printf("[%v] i: %v\n", time.Now().Format("2006-01-02 15:04:05"), i)
		time.Sleep(time.Second)
	}
}
