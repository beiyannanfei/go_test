package main

import (
	"fmt"
	"time"
)

func main() {
	c := time.Tick(time.Second * 5)
	for {
		<-c
		go f()
	}
}

func f() {
	fmt.Printf("[%v] 每5s输出一次.\n", time.Now().Format("2006-01-02 15:04:05"))
}
