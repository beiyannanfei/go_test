package main

import "fmt"

func main() {
	f := func(x, y int) int {
		return x + y
	}

	fmt.Println(f(10, 20))	//=> 30
}
