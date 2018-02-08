//Go 语言中也可以使用相互调用的递归函数：多个函数之间相互调用形成闭环。
// 因为 Go 语言编译器的特殊性，这些函数的声明顺序可以是任意的。
// 下面这个简单的例子展示了函数 odd 和 even 之间的相互调用
package main

import (
	"fmt"
)

func main() {
	fmt.Printf("%d is even: is %t\n", 16, even(16)) // 16 is even: is true
	//fmt.Printf("%d is odd: is %t\n", 17, odd(17))
	// 17 is odd: is true
	//fmt.Printf("%d is odd: is %t\n", 18, odd(18))
	// 18 is odd: is false
}

func even(nr int) bool {
	fmt.Println("even: ", nr)
	if nr == 0 {
		return true
	}
	return odd(RevSign(nr) - 1)
}

func odd(nr int) bool {
	fmt.Println("odd: ", nr)
	if nr == 0 {
		return false
	}
	return even(RevSign(nr) - 1)
}

func RevSign(nr int) int {
	if nr < 0 {
		return -nr
	}
	return nr
}

