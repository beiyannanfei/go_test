package main

import "fmt"

type people29 interface {
	speak()
}
type student29 struct {
	name string
	age  int
}

func (stu *student29) speak() {
	fmt.Println("I am a studeng, I am ", stu.age)
}

func main() {
	var p people29
	//p = student29{name: "RyuGou", age: 25}	//cannot use student29 literal (type student29) as type people29 in assignment:
	p = &student29{name: "RyuGou", age: 25}
	p.speak()
}
