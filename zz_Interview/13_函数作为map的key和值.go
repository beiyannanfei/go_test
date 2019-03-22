package main

import "fmt"

// 13~23 例子来源 https://juejin.im/post/5b5bd2ddf265da0f716c2fea?utm_source=gold_browser_extension

//Golang中函数被看做是值,函数值不可以比较，也不可以作为map的key

func main() {
	array := make(map[int]func() int)

	array[func() int { return 10 }()] = func() int {	//todo 函数值类型不能作为map的key
		return 12
	}

	fmt.Println(array)

	v := array[10]
	fmt.Println(v())

	/*
	array := make(map[func ()int]int)
	array[func()int{return 12}] = 10
	fmt.Println(array)
	不能编译通过
	在Go语言中，函数被看做是第一类值：(first-class values)：函数和其他值一样，可以被赋值，可以传递给函数，可以从函数返回。也可以被当做是一种“函数类型”。
		例如：有函数``func square(n int) int { return n * n }``，那么就可以赋值``f := square``,而且还可以``fmt.Println(f(3))``（将打印出“9”）。
	Go语言函数有两点很特别：
		+ 函数值类型不能作为map的key
		+ 函数值之间不可以比较，函数值只可以和nil作比较，函数类型的零值是``nil``
	*/
}
