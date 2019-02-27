package main

import "log"

//https://studygolang.com/articles/17581

func main() {
	bb1()
	bb2()
}



func bb2() {
	log.Printf("===============bb2 begin================\n\n")
	//闭包是可以包含自由变量的代码块，这些变量不在这个代码块内或者任何全局上下文中定义，而是在定义代码块的环境中定义。
	//闭包的价值在于可以作为函数对象或者匿名函数，存储到变量中作为参数传递给其他函数，能够被函数动态创建和返回。

	var j int = 5
	a := func() func() {
		var i int = 10
		return func() {
			log.Printf("i: %d, j: %d\n\n", i, j)
		}
	}()

	a()
	//i: 10, j: 5
	j *= 2
	//i: 10, j: 10
	a()

	log.Printf("===============bb2 end================\n\n")
}

func bb1() {
	log.Printf("===============bb1 begin================\n\n")

	f := func(x, y int) int {
		return x + y
	}
	log.Printf("f(4, 5) = %v\n", f(4, 5))

	z := func(x, y int) int {
		log.Printf("匿名函数，直接执行, x: %v, y: %v\n", x, y)
		return x + y
	}(6, 7)

	//匿名函数，直接执行, x: 6, y: 7
	log.Printf("z = %v\n\n", z)
	//z = 13

	log.Printf("===============bb1 end================\n\n")
}
