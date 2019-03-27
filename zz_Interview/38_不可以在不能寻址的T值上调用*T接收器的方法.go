package main

import "fmt"

type Lili38 struct {
	Name string
}

func (Lili *Lili38) fmtPointer() {
	fmt.Println("pointer")
}

func (Lili Lili38) fmtReference() {
	fmt.Println("reference")
}

func main() {
	li := Lili38{}
	li.fmtPointer() //pointer

	Lili38{}.fmtReference()	//reference
	Lili38{}.fmtPointer()	//cannot call pointer method on Lili38 literal
	//todo 假设T类型的方法上接收器既有T类型的，又有*T指针类型的，那么就不可以在不能寻址的T值上调用*T接收器的方法 main主函数中的“li”是一个变量，li的虽然是类型Lili，但是li是可以寻址的，&li的类型是*Lili，因此可以调用*Lili的方法。
}
