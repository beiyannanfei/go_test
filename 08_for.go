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
}

func for2() {
	var a int
	var b = 6
	for a < b {
		a++;
		fmt.Printf("a is %d\n", a)
	}

}

func for1() {
	for a := 0; a < 5; a++ {
		fmt.Printf("a is %d\n", a)
	}
}
