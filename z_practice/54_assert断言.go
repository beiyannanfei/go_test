package main

import (
	"fmt"
	"reflect"
)

//Go语言里面有一个语法，可以直接判断是否是该类型的变量: value, ok = element.(T)，这里value就是变 量的值，ok是一个bool类型，element是interface变量，T是断言的类型。
//如果element里面确实存储了T类型的数值，那么ok返回true，否则返回false。

type Element interface{}
type List []Element

type Person struct {
	name string
	age  int
}

func (p *Person) String() string {
	return fmt.Sprintf("(name: %v, age: %v years)", p.name, p.age)
}

func main() {
	list := make(List, 4)
	list[0] = 1;
	list[1] = "Hello"
	list[2] = Person{"Dennis", 70}
	list[3] = 1.23

	for index, element := range list {
		if value, ok := element.(int); ok {
			fmt.Printf("list[%v] is an int and its value is %v\n", index, value)
			continue
		}

		if value, ok := element.(string); ok {
			fmt.Printf("list[%v] is an string and its value is %v\n", index, value)
			continue
		}

		if value, ok := element.(Person); ok {
			fmt.Printf("list[%v] is an Person and its value is %v\n", index, value)
			continue
		}

		fmt.Printf("list[%v] is of a different type: %v\n", index, reflect.TypeOf(element))
	}
}
