package main

import (
	"fmt"
	"time"
	"reflect"
)

//https://www.cnblogs.com/suoning/p/7237444.html
//应用场景，如果某个goroutine panic了，而且这个goroutine里面没有捕获(recover)，那么整个进程就会挂掉。所以，好的习惯是每当go产生一个goroutine，就需要写下recover

var domainSyncChan = make(chan int, 10)

func domainPut(num int) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("error to chan put.", err, reflect.TypeOf(err))
		}
	}()
	domainSyncChan <- num

	panic(fmt.Sprintf("error....%v", num))
}

func main() {
	for i := 0; i < 10; i++ {
		go domainPut(i)
	}

	time.Sleep(time.Second * 2)
}
