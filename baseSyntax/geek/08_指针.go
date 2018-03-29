package main

import "fmt"

func main() {
	var i1 = 5
	fmt.Printf("An integer: %d, it's location in memory: %p\n", i1, &i1)

	var intP *int
	intP = &i1
	fmt.Printf("The value at memory location %p is %d\n", intP, *intP)

	strPointer()

	nilPointer()
}

func strPointer() {
	s := "good bye"
	var p *string = &s
	*p = "ciao"
	fmt.Printf("Here is the pointer p: %p\n", p)  //Here is the pointer p: 0xc4200701d0
	fmt.Printf("Here is the string *p: %s\n", *p) //Here is the string *p: ciao
	fmt.Printf("Here is the string s: %s\n", s)   //Here is the string s: ciao
}

func nilPointer() {
	//var p *int = nil //对一个空指针的反向引用是不合法的
	//*p = 0
}
