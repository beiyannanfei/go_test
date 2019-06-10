package main

import "fmt"

type student45 struct {
	Name string
	Age  int
}

func main() {
	m := make(map[string]*student45)
	stus := []student45{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 25},
		{Name: "wang", Age: 26},
	}

	for _, stu := range stus { //todo range 循环，会重用地址
		m[stu.Name] = &stu
		/* 正确写法
		temp := stu
		m[stu.Name] = &temp
		*/
	}

	for key, value := range m {
		fmt.Printf("key: %s, value: %v\n", key, value)
	}
	/*
	key: zhou, value: &{wang 26}
	key: li, value: &{wang 26}
	key: wang, value: &{wang 26}
	*/
}
