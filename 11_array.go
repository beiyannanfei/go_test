package main

import "fmt"

func main() {
	dingyi()
}

func dingyi() { //数组的定义
	//初始化数组中 {} 中的元素个数不能大于 [] 中的数字
	var balance = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}

	//如果忽略 [] 中的数字不设置数组大小，Go 语言会根据元素的个数来设置数组的大小
	var balance1 = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}

	for i := 0; i < len(balance); i++ {
		fmt.Printf("balance[%d] = %f\n", i, balance[i])
	}
	fmt.Printf("=========== all balance is %#v\n", balance)
	for i := 0; i < len(balance1); i++ {
		fmt.Printf("balance1[%d] = %f\n", i, balance1[i])
	}
	fmt.Printf("=========== all balance1 is %v\n", balance1)
}
