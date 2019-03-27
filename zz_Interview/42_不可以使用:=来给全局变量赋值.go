package main

import (
	"errors"
	"fmt"
)

var varTest string

func test42() {
	varTest, err := function42()
	fmt.Println(err.Error())
	/* 正确写法
	err := errors.New("error")
	varTest, err = function()
	fmt.Println(err.Error())
	*/
}

func function42() (string, error) {
	return "hello world", errors.New("error")
}

func main() {
	test42()	//varTest declared and not used
	//todo 在test方法中，如果使用varTest, err := function()这种方式的话，相当于在函数中又定义了一个和全局变量varTest名字相同的局部变量，而这个局部变量又没有使用，所以会编译不通过。
}
