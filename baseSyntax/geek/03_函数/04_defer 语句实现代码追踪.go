package main

import "fmt"

func trace(s string)   { fmt.Println("entering:", s) }
func untrace(s string) { fmt.Println("leaving:", s) }

func a() {
	trace("a")          //step 3
	defer untrace("a")  //step 5
	fmt.Println("in a") //step 4
}

func b() {
	trace("b")          //step 1
	defer untrace("b")  //step 6
	fmt.Println("in b") //step 2
	a()
}

func main() {
	b()
}
