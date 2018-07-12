package main

import "fmt"

func main() {
	var a = 10;
LOOP:
	for a < 20 {
		if a == 15 {
			a++
			goto LOOP
		}
		fmt.Printf("a is %d\n", a)
		a++
	}
	/*
	a is 10
	a is 11
	a is 12
	a is 13
	a is 14
	a is 16
	a is 17
	a is 18
	a is 19
	*/
}
