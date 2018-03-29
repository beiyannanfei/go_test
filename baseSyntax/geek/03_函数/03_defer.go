//关键字 defer 允许我们进行一些函数执行完成后的收尾工作

package main

import "fmt"

func main() {
	defer1()
	defer2()
	defer3()
}

func defer3() {
	for i := 0; i < 5; i++ {
		defer fmt.Printf("%d ", i) //当有多个 defer 行为被注册时，它们会以逆序执行（类似栈，即后进先出）out: 43210
	}
}

func defer2() {
	i := 0
	defer fmt.Println(i) //0
	i++
	return
}

func defer1() {
	defer function2() //最后才会执行
	fmt.Printf("11111\n")
	fmt.Printf("22222\n")
	return
}

func function2() {
	fmt.Printf("33333\n")
}
