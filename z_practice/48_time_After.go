package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("----------------- After -----------------")
	//After 会在另一线程经过时间段 d 后向返回值发送当时的时间。等价于NewTimer(d).C。
	fmt.Println("start:", time.Now().Format("2006-01-02 15:04:05"))

	timer := time.After(time.Second * 2)

	select {
	case t := <-timer:
		fmt.Println("get timer:", t.Format("2006-01-02 15:04:05"))
	}

	fmt.Println("stopped:", time.Now().Format("2006-01-02 15:04:05"))
}
