/**
 * Created by wyq on 18/2/3.
 */
package main

import (
	"fmt"
	"unsafe"
)

func constExp() {
	const (
		a = "abc"
		b = len(a)
		c = unsafe.Sizeof(a)
	)
	println(a, b, c) //abc 3 16
}

func emit() { //常量枚举
	const (
		Unknown = 0
		Female  = 1
		Male    = 2
	)
	println(Unknown, Female, Male) //0 1 2
}

//iota，特殊常量，可以认为是一个可以被编译器修改的常量。
//在每一个const关键字出现时，被重置为0，然后再下一个const出现之前，每出现一次iota，其所代表的数字会自动增加1。
func iota1() {
	const (
		a = iota
		b = iota
		c = iota
	)
	println("iota1: ", a, b, c)
}

//第一个 iota 等于 0，每当 iota 在新的一行被使用时，它的值都会自动加 1；所以 a=0, b=1, c=2 可以简写为如下形式：
func iota2() {
	const (
		a = iota
		b
		c
	)
	println("iota2: ", a, b, c)
}

func iota3() {
	const (
		a = iota //0
		b        //1
		c        //2
		d = "ha" //独立值ha，iota += 1
		e        //"ha"   iota += 1
		f = 100  //100 iota +=1
		g        //100  iota +=1
		h = iota //7,恢复计数
		i        //8
	)
	println("iota3: ", a, b, c, d, e, f, g, h, i) //iota3:  0 1 2 ha ha 100 100 7 8
}

func iota4() {
	const (
		i = 1 << iota	// <=> 1 << 0
		j = 3 << iota	// <=> 3 << 1 <=> 110 = 6
		k				// <=> 3 << 2 <=> 1100 = 12
		l				// <=> 3 << 3 <=> 11000 = 24
	)
	println("i=", i)
	println("j=", j)
	println("k=", k)
	println("l=", l)
}

func main() {
	const LENGTH int = 10
	const WIDTH int = 5
	var area int
	const a, b, c = 1, false, "str" //多重赋值

	area = LENGTH * WIDTH
	fmt.Printf("面积为 : %d", area)
	println()
	println(a, b, c) //1 false str
	emit()
	constExp()
	iota1()
	iota2()
	iota3()
	iota4()
}
