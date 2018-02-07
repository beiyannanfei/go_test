package main

import "fmt"

func main() {
	fmt.Printf(`This is a raw string \n`) // \n会原样输出
	fmt.Println()
	str := "Beginning of the string " + //+ 必须放在第一行
		"second part of the string"
	fmt.Println(str)
	s := "hel" + "lo,"
	s += "world!"
	fmt.Println(s) //输出 “hello, world!”
}
