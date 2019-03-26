package main

import "fmt"

func main() {
	a := 1
	defer print28(function(a)) //1
	a = 2
}

func function(num int) int {
	return num
}

func print28(num int) {
	fmt.Println(num)
}
