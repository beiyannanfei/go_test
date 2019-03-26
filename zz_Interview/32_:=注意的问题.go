package main

import "fmt"

a := 10 // todo 使用:=定义的变量，仅能使用在函数内部。

func main() {
	b := 20
	fmt.Println(a, b)
}

//todo 在定义多个变量的时候:=周围不一定是全部都是刚刚声明的，有些可能只是赋值，例如下面的err变量
//in, err := os.Open(infile)
//out, err := os.Create(outfile)
