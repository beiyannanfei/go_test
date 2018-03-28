package main

import (
	"fmt"
)

const c = "C"

var v int = 5

type T struct{}

func init() { // initialization of package
	fmt.Println("... init fun ...")
}

func main() {
	fmt.Println("... main ...")
	var a int
	Func1()
	// ...
	fmt.Println(a)
	var t = new(T)
	t.Method1()
}

func (t T) Method1() {
	//...
	fmt.Println("... Method1 ...")
}

func Func1() { // exported function Func1
	//...
	fmt.Println("... Func1 ...")
}
