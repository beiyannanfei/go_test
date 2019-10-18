package main

import "fmt"

func s1() {
	s := make([]int, 5)		//创建一个长度为5的切片，并用对应类型的默认值初始化
	fmt.Printf("%p\n", s)
	s = append(s, 1, 2, 3)
	fmt.Printf("%p\n", s)
	fmt.Println(s)
}

func s2() {
	s := make([]int, 0, 5)
	fmt.Printf("%p\n", s)
	s = append(s, 1, 2, 3)
	fmt.Printf("%p\n", s)
	fmt.Println(s)
}

func main() {
	s1()
	fmt.Println("-----------------------------")
	s2()
}

/*
0xc000098000
0xc0000a4000
[0 0 0 0 0 1 2 3]
-----------------------------
0xc000098030
0xc000098030
[1 2 3]
*/
