package main

import "fmt"

func main() {
	x := []int{2, 3, 5, 7, 11}
	y := x[1:3]
	fmt.Println(x, y)
	fmt.Println(len(x), cap(x))
	fmt.Println(len(y), cap(y))

}
