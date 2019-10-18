package main

import "fmt"

type student struct {
	Name string
	age  int
}

func pase_student() map[string]*student {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", age: 24},
		{Name: "li", age: 23},
		{Name: "wang", age: 22},
	}

	//stu变量的地址始终保持不变，每次遍历仅进行struct值拷贝，故m[stu.Name]=&stu实际上一直指向同一个地址，最终该地址的值为遍历的最后一个struct的值拷贝
	for _, stu := range stus {
		fmt.Printf("%v \t %p\n", stu, &stu)
		m[stu.Name] = &stu
	}

	return m
}

func pase_student1() map[string]*student {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", age: 24},
		{Name: "li", age: 23},
		{Name: "wang", age: 22},
	}

	for i, _ := range stus {
		fmt.Printf("%v \t %p\n", stus[i], &stus[i])
		m[stus[i].Name] = &stus[i]
	}

	return m
}

func main() {
	students := pase_student()
	for k, v := range students {
		fmt.Printf("key = %s, value = %v\n", k, v)
	}
	fmt.Println("-------------------------------------------")

	students1 := pase_student1()
	for k, v := range students1 {
		fmt.Printf("key = %s, value = %v\n", k, v)
	}
}

/*
{zhou 24} 	 0xc00000c060
{li 23} 	 0xc00000c060
{wang 22} 	 0xc00000c060
key = zhou, value = &{wang 22}
key = li, value = &{wang 22}
key = wang, value = &{wang 22}
-------------------------------------------
{zhou 24} 	 0xc0000a6000
{li 23} 	 0xc0000a6018
{wang 22} 	 0xc0000a6030
key = zhou, value = &{zhou 24}
key = li, value = &{li 23}
key = wang, value = &{wang 22}
*/
