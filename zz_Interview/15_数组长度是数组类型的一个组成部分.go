package main

import (
	"fmt"
	"reflect"
)

// 13~23 例子来源 https://juejin.im/post/5b5bd2ddf265da0f716c2fea?utm_source=gold_browser_extension

func main() {
	arrayA := [...]int{1, 2, 3}
	arrayB := [...]int{1, 2, 3, 4}
	fmt.Printf("arrayA type: %s\n", reflect.TypeOf(arrayA))
	fmt.Printf("arrayB type: %s\n", reflect.TypeOf(arrayB))
	fmt.Println(reflect.TypeOf(arrayA) == reflect.TypeOf(arrayB)) //todo 数组长度是数组类型的一个组成部分，因此[3]int和[4]int是两种不同的数组类型
}
