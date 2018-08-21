package main

import (
	"fmt"
	"time"
)

//如果map由多协程同时读和写就会出现 fatal error:concurrent map read and map write的错误
import (
	"strconv"
)

func main() {
	c := make(map[string]int)

	go func() { //开一个协程写map
		for i := 0; i < 100000; i++ {
			c[strconv.Itoa(i)] = i
		}
	}()

	go func() { //开一个协程读map
		for i := 0; i < 100000; i++ {
			fmt.Println(c[strconv.Itoa(i)])
		}
	}()

	time.Sleep(time.Second * 20)
}
//解决方案见下例