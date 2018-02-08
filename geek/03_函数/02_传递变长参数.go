package main

import "fmt"

func main() {
	x := min(1, 3, 2, 0)
	fmt.Printf("The minimum is: %d\n", x)
	arr := []int{7, 9, 3, 5, 1}
	x = min(arr...) //如果参数被存储在一个数组 arr 中，则可以通过 arr... 的形式来传递参数调用变参函数。
	fmt.Printf("The minimum in the array arr is: %d", x)
}

func F1(s ...string) {
	F2(s...) //一个接受变长参数的函数可以将这个参数作为其它函数的参数进行传递
	F3(s)
}

func F2(s ...string) {}
func F3(s []string)  {}

func min(a ...int) int { //a为参数的列表切片
	if len(a) == 0 {
		return 0
	}
	min := a[0]
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}
