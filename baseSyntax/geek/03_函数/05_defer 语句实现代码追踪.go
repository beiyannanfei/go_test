package main

import "fmt"

func trace1(s string) string {
	fmt.Println("entering:", s)
	return s
}

func un(s string) {
	fmt.Println("leaving:", s)
}

func a1() {
	defer un(trace1("a1"))
	fmt.Println("in a1")
}

func b1() {
	defer un(trace1("b1"))
	fmt.Println("in b1")
	a1()
}

func main() {
	b1()
}

/*
entering: b1
in b1
entering: a1
in a1
leaving: a1
leaving: b1
*/