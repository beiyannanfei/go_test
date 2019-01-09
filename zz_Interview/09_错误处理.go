package main

import "log"

//https://studygolang.com/articles/17581

//3个重要关键字defer、panic、recover
//defer是函数结束后执行，呈先进后出；
//panic是程序出现无法修复的错误时使用，但会让defer执行完；
//recover会修复错误，不至于程序终止。当不确定函数不会出错时使用defer+recover
func fixError() {
	if r := recover(); r != nil {
		log.Printf("err caught: %v\n", r)
	} else {
		log.Println("no error")
	}
}

func myDivide(x, y int) int {
	return x / y
}

func testFun() {
	defer fixError()
	myDivide(6, 0)
}

func main() {
	testFun()
	log.Println("main end")
	log.Println("=========================================")


}

//err caught: runtime error: integer divide by zero
//main end
