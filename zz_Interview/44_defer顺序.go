package main

import "fmt"

func main() {
	defer_call44()
}

func defer_call44() {
	defer func() {
		fmt.Println("打印前")
	}()

	defer func() {
		fmt.Println("打印中")
	}()

	defer func() {
		fmt.Println("打印后")
	}()

	panic("触发异常")
}
/*
打印后
打印中
打印前
*/