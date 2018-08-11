package main

import (
	"fmt"
	"reflect"
)

type X1 struct{}

func (X1) SumTest(x int, y int) (int, error) {
	return x + y, fmt.Errorf("err: %v", x+y)
}

func methodCall() {
	var a X1
	v := reflect.ValueOf(a)
	m := v.MethodByName("SumTest")
	in := []reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf(2),
	}
	out := m.Call(in)
	for i, v := range out {
		fmt.Printf("i: %v, v: %v\n", i, v)
	}
	/*
	i: 0, v: 3
	i: 1, v: err: 3
	*/
}

func (X1) FormatTest(s string, a ...interface{}) string {
	return fmt.Sprintf(s, a...)
}

func methodCallSlice() {
	var a X1
	v := reflect.ValueOf(a)
	m := v.MethodByName("FormatTest")
	out := m.Call([]reflect.Value{
		reflect.ValueOf("%s = %d"), // 所有参数都须处理
		reflect.ValueOf("x"),
		reflect.ValueOf(100),
	})
	fmt.Println("call out:", out) //call out: [x = 100]

	out = m.CallSlice([]reflect.Value{
		reflect.ValueOf("%v = %v"),
		reflect.ValueOf([]interface{}{"xyx", 789}),
	})
	fmt.Println("callslice out:", out) //callslice out: [xyx = 789]
}

func main() {
	fmt.Println("------ 动态调用方法，只须按 In 列表准备好所需参数即可 ------")
	methodCall()

	fmt.Println("------ 对于变参来说，用 CallSlice() 要更方便一些 ------")
	methodCallSlice()
}
