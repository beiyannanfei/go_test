package main

import "fmt"

//https://i6448038.github.io/2017/07/28/GolangDetails/

//todo 多个defer出现的时候，多个defer之间按照LIFO（后进先出）的顺序执行
func main() {
	defer func() {
		fmt.Println(1)
	}()

	defer func() {
		fmt.Println(2)
	}()

	defer func() {
		fmt.Println(3)
	}()
}
