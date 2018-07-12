package main

import "fmt"

func main() {
	var i = 65
	fmt.Printf("%d 的阶乘是 %d\n", i, Factorial(uint64(i)))
	fmt.Println("---------------------------")
	for i := 0; i < 65; i++ {
		fmt.Printf("%d-%d\n", i, fibonacci(i))
	}
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-2) + fibonacci(n-1)
}

func Factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}
