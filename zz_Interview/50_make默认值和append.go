package main

import "fmt"

func main() {
	s := make([]int, 5)
	s = append(s, 1, 2, 3)
	fmt.Println(s) //[0 0 0 0 0 1 2 3]
	fmt.Println("------------------------")
	s1 := make([]int, 0)
	s1 = append(s1, 1, 2, 3)
	fmt.Println(s1) //[1 2 3]
}

//todo make初始化是由默认值的哦，此处默认值为0
