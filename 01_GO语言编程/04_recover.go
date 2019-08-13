package main

import (
	"log"
	"fmt"
)

func foo() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Runtime error caught: %v", r)
		}
	}()

	a := 1
	b := 2
	c := 1
	d := b / (a - c)
	fmt.Println(d)
}

func main() {
	foo()
	for i := 0; i < 10; i++ {
		fmt.Println("hello world")
	}

	var a int = 10
	var b = 10
	c := 20
	fmt.Println(a, b, c)
}
