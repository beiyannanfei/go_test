package main

import "fmt"

func v1() {
	var msgs []func()
	array := []string{
		"10", "20", "30", "40",
	}

	for _, e := range array {
		msgs = append(msgs, func() {
			fmt.Println(e)
		})
	}

	for _, v := range msgs {
		v()
	}
	/*
		40
		40
		40
		40
	*/
}

func v2()  {
	var msgs []func()
	array := []string{
		"10", "20", "30", "40",
	}

	for _, e := range array {
		elem := e	//todo 其实就加了条elem := e看似多余，其实不，这样一来，每次循环后每个匿名函数中保存的就都是当时局部变量elem的值，这样的局部变量定义了4个，每次循环生成一个。
		msgs = append(msgs, func() {
			fmt.Println(elem)
		})
	}

	for _, v := range msgs {
		v()
	}
	/*
		10
		20
		30
		40
	*/
}

func main() {
	v1()
	fmt.Println("================================")
	v2()
}
