package main

import "fmt"

func main() {
	goto1()
	goto2()
	println()

}

func goto3() {
	/*a := 1
	goto TARGET // compile error		//goto TARGET jumps over declaration of b
	b := 9
TARGET:
	b += a
	fmt.Printf("a is %v *** b is %v", a, b)*/
}

func goto2() {
	i := 0
HERE:
	fmt.Print(i)
	i++
	if i == 5 {
		return
	}
	goto HERE
}

func goto1() {
LABEL1:
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5; j++ {
			if j == 4 {
				continue LABEL1
			}
			fmt.Printf("i is: %d, and j is: %d\n", i, j)
		}
	}

}
