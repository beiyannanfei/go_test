package main

import "fmt"

type littleGirl22 struct {
	Name string
	Age  int
}

func (this *littleGirl22) changeName(name string) {
	fmt.Println(name)
}

func main() {
	little := &littleGirl22{Name: "Rose", Age: 1}

	//todo Go语言中，允许方法用nil指针作为其接收器，也允许函数将nil指针作为参数。而上述代码中的littleGirl不是指针类型，改为*littleGirl，然后变量little赋值为&littleGirl{Name:"Rose", Age:1}就可以编译通过了。	并且，nil对于对象来说是合法的零值的时候，比如map或者slice，也可以编译通过并正常运行。

	little = nil
	little.changeName("yoyo")
	fmt.Println(little)
}
