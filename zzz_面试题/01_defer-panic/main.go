package main

import "fmt"

func main() {
	defer_call()
}

func defer_call() {
	defer func() {
		fmt.Println("before log")
	}()

	defer func() {
		fmt.Println("logging...")
	}()

	defer func() {
		fmt.Println("after log")
	}()

	panic("exception")
}
/*
after log
panic: exception	//异常触发位置不确定
logging...
before log
*/