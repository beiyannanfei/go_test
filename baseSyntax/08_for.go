package main

import "fmt"

func main() {
	for1()
	fmt.Println("-------------------")
	for2()
	fmt.Println("-------------------")
	for3()
}

func for3() {
	numbers := [6]int{1, 2, 3, 5}
	for i, x := range numbers {
		fmt.Printf("index %d value is %d\n", i, x)
	}
	/*
	index 0 value is 1
	index 1 value is 2
	index 2 value is 3
	index 3 value is 5
	index 4 value is 0
	index 5 value is 0
	*/
}

func for2() {
	var a int
	var b = 6
	for a < b {
		a++
		fmt.Printf("a is %d\n", a)	//1 2 3 4 5 6
	}

}

func for1() {
	for a := 0; a < 5; a++ {
		fmt.Printf("a is %d\n", a)	//0 1 2 3 4
	}
}
