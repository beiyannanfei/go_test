//与字符串相关的类型转换都是通过 strconv 包实现的。
//该包包含了一些变量用于获取程序运行的操作系统平台下 int 类型所占的位数，如：strconv.IntSize。
//任何类型 T 转换为字符串总是成功的。
//针对从数字类型转换到字符串，Go 提供了以下函数：
//strconv.Itoa(i int) string 返回数字 i 所表示的字符串类型的十进制数。
//strconv.FormatFloat(f float64, fmt byte, prec int, bitSize int) string
// 将 64 位浮点型的数字转换为字符串，其中 fmt 表示格式（其值可以是 'b'、'e'、'f' 或 'g'），
// prec 表示精度，bitSize 则使用 32 表示 float32，用 64 表示 float64。
//将字符串转换为其它类型 tp 并不总是可能的，可能会在运行时抛出错误 parsing "…": invalid argument。
//针对从字符串类型转换为数字类型，Go 提供了以下函数：
//strconv.Atoi(s string) (i int, err error) 将字符串转换为 int 型。
//strconv.ParseFloat(s string, bitSize int) (f float64, err error) 将字符串转换为 float64 型。
//利用多返回值的特性，这些函数会返回 2 个值，第 1 个是转换后的结果（如果转换成功），
// 第 2 个是可能出现的错误，因此，我们一般使用以下形式来进行从字符串到其它类型的转换：
//val, err = strconv.Atoi(s)
package main

import (
	"fmt"
	"strconv"
)

func main() {
	var orig string = "666"
	var an int
	var newS string

	fmt.Printf("The size of ints is: %d\n", strconv.IntSize)		//The size of ints is: 64

	an, _ = strconv.Atoi(orig)
	fmt.Printf("The integer is: %d\n", an)						//The integer is: 666
	an = an + 5
	newS = strconv.Itoa(an)
	fmt.Printf("The new string is: %s\n", newS)					//The new string is: 671

	orig = "123.456"
	var f float64
	f, _ = strconv.ParseFloat(orig, 64)
	fmt.Printf("The float is: %v\n", f)							//The float is: 123.456
	f += 0.123456
	fmt.Printf("The float is: %v\n", f)							//The float is: 123.57945600000001
	newS = strconv.FormatFloat(f, 'f', 7, 64)
	fmt.Printf("The new string is: %s\n", newS)					//The new string is: 123.5794560
}
