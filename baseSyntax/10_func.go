package main

import (
	"fmt"
	"math"
)

func main() {
	res := max(10, 20)
	fmt.Printf("the max num is %d\n", res)		//the max num is 20
	a, b := swap("hello", "world")
	fmt.Println(a, b)	//world hello
	var c, d = 100, 200
	fmt.Printf("before swap c: %d, d: %d\n", c, d)	//before swap c: 100, d: 200
	quot(&c, &d)
	fmt.Printf("after swap c: %d, d: %d\n", c, d)		//after swap c: 200, d: 100
	funAsVal()		//3
	fmt.Println("=========================")
	nextNumber := closureFun()
	fmt.Println(nextNumber()) //1
	fmt.Println(nextNumber()) //2
	fmt.Println(nextNumber()) //3

	nextNumber1 := closureFun()
	fmt.Println(nextNumber1()) //1
	fmt.Println(nextNumber1()) //2
	fmt.Println("=========================")

	add_func := add(1, 2)
	fmt.Println(add_func()) //1 3
	fmt.Println(add_func()) //2 3
	fmt.Println(add_func()) //3 3

	fmt.Println("=========================")
	var c1 Circle
	c1.radius = 10.00
	fmt.Println("area of Circle(c1) = ", c1.getArea()) //area of Circle(c1) =  314
}

type Circle struct {
	radius float64
}

func (c Circle) getArea() float64 {
	return c.radius * c.radius * 3.14
}

// 闭包使用方法-带参数的闭包函数调用
func add(x1, x2 int) func() (int, int) {
	i := 0
	return func() (int, int) {
		i++
		return i, x1 + x2
	}
}

func closureFun() func() int { //闭包函数
	i := 0
	return func() int {
		i += 1
		return i
	}
}

func funAsVal() { //函数作为值
	getSquareRoot := func(x float64) float64 {
		return math.Sqrt(x)
	}
	fmt.Println(getSquareRoot(9)) //3
}

func quot(x *int, y *int) { //引用传值
	var temp = *x
	*x = *y
	*y = temp
}

func swap(a, b string) (string, string) { //返回多个值
	return b, a
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}
