/**
 * Created by wyq on 18/2/3.
 */

package main

var a = "菜鸟教程"
var b string = "runoob.com"
var c bool
//d := 10  这种格式的变量声明只能在函数中

var x, y int //多变量声明
var ( // 这种因式分解关键字的写法一般用于声明全局变量
	a1 int
	b1 bool
)

var c1, d1 int = 1, 2
var e1, f1 = 123, "hello"

func main() {
	d := 123                                                        //这种格式的变量声明只能在函数中
	println(a, b, c, d, "\n=================")                      //out: 菜鸟教程 runoob.com false 123
	println(x, y, a1, b1, c1, d1, e1, f1, "\n--------------------") //0 0 0 false 1 2 123 hello
	//g1 := "abc" //err: ./02_Go 语言变量.go:25: g1 declared and not used 在函数中声明的变量不使用会报错

}
