package main

import "fmt"

//https://i6448038.github.io/2018/07/18/golang-mistakes/

type Father43 interface {
	Hello()
}

type Child43 struct {
	Name string
}

func (s Child43) Hello() {

}

func f43(out *Father43) {
	if out != nil {
		fmt.Println("surprise!")
	}
}

func main() {
	var buf Child43 //正确写法 var buf  Father43
	buf = Child43{}
	f43(&buf) //todo 接口类型的变量可以被赋值为实现接口的结构体的实例，但是并不能代表接口的指针可以被赋值为实现接口的结构体的指针实例
}
