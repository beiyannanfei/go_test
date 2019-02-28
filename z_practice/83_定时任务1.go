package main

import (
	"time"
	"fmt"
)

func main() {
	for range time.Tick(time.Second * 5) {
		fmt.Printf("[%v] 每5s输出一次.\n", time.Now().Format("2006-01-02 15:04:05"))
	}
}
