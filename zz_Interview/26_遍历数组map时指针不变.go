package main

import "fmt"

type student26 struct {
	Name string
	Age  int
}

func main() {
	var stus []student26

	stus = []student26{
		{Name: "one", Age: 18},
		{Name: "two", Age: 19},
	}

	data := make(map[int]*student26)

	//todo 用for range来遍历数组或者map的时候，被遍历的指针是不变的，每次遍历仅执行struct值的拷贝
	for i, v := range stus {
		data[i] = &v
	}

	for i, v := range data {
		fmt.Printf("key=%d, value=%v\n", i, v)
	}

	//key=0, value=&{two 19}
	//key=1, value=&{two 19}
}
