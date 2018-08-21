package main

import (
	"fmt"
	"reflect"
	"time"
)
//https://blog.csdn.net/skh2015java/article/details/60330975
func xrange() chan int { //xrange 用来生成自增的整数
	var ch = make(chan int) //无缓冲信道

	go func() { //开出一个goroutine
		for i := 0; ; i++ {
			ch <- i //因为信道无缓冲，所以存入一个数据后如果不被取走，会一直阻塞
			fmt.Println("xrange store i:", i)
		}
	}()

	return ch
}

func main() {
	generage := xrange()
	fmt.Println("generate type: ", reflect.TypeOf(generage))

	for i := 0; i < 10; i++ {
		v, ok := <-generage
		fmt.Printf("v: %v, ok: %v\n", v, ok)
		time.Sleep(time.Second)
	}
}
