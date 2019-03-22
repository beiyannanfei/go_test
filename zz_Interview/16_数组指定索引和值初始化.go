package main

import "fmt"

func main() {
	arrayA := [...]int{0: 1, 2: 1, 3: 4}
	fmt.Println(arrayA)			//[1 0 1 4]
	fmt.Println(len(arrayA))	//4
	//todo 定义了一个数组长度为4的数组，指定索引的数组长度和最后一个索引的数值相关，例如:r := [...]int{99:-1}就定义了一个含有100个元素的数组r，最后一个元素输出化为-1，其他的元素都是用0初始化。
}
