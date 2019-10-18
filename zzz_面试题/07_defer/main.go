package main

import "fmt"

func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func main() {
	a := 1
	b := 2
	//defer是在函数末尾的return前执行，先进后执行
	//函数调用时 int 参数发生值拷贝
	defer calc("1", a, calc("10", a, b)) //首先执行calc("10", 1, 2), 然后执行defer("1", 1, 3)
	a = 0
	defer calc("2", a, calc("20", a, b)) //首先执行calc("20", 0, 2), 然后执行defer("2", 0, 2)
	b = 1
}

/*
10 1 2 3
20 0 2 2
2 0 2 2
1 1 3 4
*/
