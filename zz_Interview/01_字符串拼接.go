package main

import "fmt"

/*
5.   【初级】关于字符串连接，下面语法正确的是（）
A. str := ‘abc’ + ‘123’
B. str := "abc" + "123"
C. str ：= '123' + "abc"
D. fmt.Sprintf("abc%d", 123)

参考答案：BD
*/

func main() {
	str1 := "abc"
	str2 := "123"
	str3 := str1 + str2
	fmt.Println("st3 = ", str3)

	str4 := fmt.Sprintf("abc%d", 123)
	fmt.Println("str4 = ", str4)
}
