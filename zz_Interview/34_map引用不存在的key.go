package main

import "fmt"

func main() {
	newMap := make(map[string]int)
	fmt.Println(newMap["a"])	//out: 0
}
//todo 不报错。不同于PHP，Golang的map和Java的HashMap类似，Java引用不存在的会返回null，而Golang会返回初始值