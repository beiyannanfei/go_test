package main

import "fmt"

func main() {
	fmt.Printf("Multiply 2 * 5 * 6 = %d\n", MultiPly3Nums(2, 5, 6))		//Multiply 2 * 5 * 6 = 60
	fmt.Println(getX2AndX3(10))	//20 30
	fmt.Println(getX2AndX3_2(100))		//200 300
	i1, _, f1 := ThreeValues()
	fmt.Printf("The int: %d, the float: %f \n", i1, f1)	//The int: 5, the float: 7.500000
	min, max := MinMax(78, 65)
	fmt.Printf("Minmium is: %d, Maximum is: %d\n", min, max)	//Minmium is: 65, Maximum is: 78

	n := 0
	reply := &n
	Multiply(10, 5, reply)
	fmt.Println("Multiply:", *reply) // Multiply: 50
}

func Multiply(a, b int, reply *int) { //引用传值
	*reply = a * b
}

func MinMax(a int, b int) (min int, max int) {
	if a < b {
		min = a
		max = b
	} else { // a = b or a < b
		min = b
		max = a
	}
	return
}

func ThreeValues() (int, int, float32) {
	return 5, 6, 7.5
}

func getX2AndX3_2(input int) (x2 int, x3 int) {
	x2 = 2 * input
	x3 = 3 * input
	return //<=> return x2, x3
}

func getX2AndX3(input int) (int, int) {
	return 2 * input, 3 * input
}

func MultiPly3Nums(a int, b int, c int) int {
	// var product int = a * b * c
	// return product
	return a * b * c
}
