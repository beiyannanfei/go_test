package main

import "fmt"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	//todo panic中可以传任何值，不仅仅可以传string
	panic([]int{12312})
}
