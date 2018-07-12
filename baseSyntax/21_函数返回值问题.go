package main

import "fmt"

func main() {
	var res1, res2 int
	res1, _ = test(1)	//按照函数定义顺序返回
	fmt.Println(res1)
	_, res2 = test(2)
	fmt.Println(res2)
	res4, res3 := test(0)
	fmt.Println(res3, res4)
}

func test(n int) (res1 int, res2 int) {
	if n == 1 {
		res1 = 10	//当返回一个值的时候直接采用赋值
		return
	} else if n == 2 {
		res2 = 20
		return
	} else {
		return 100, 200
	}
}
