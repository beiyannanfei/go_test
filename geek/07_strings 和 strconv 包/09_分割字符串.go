//strings.Fields(s) 将会利用 1 个或多个空白符号来作为动态长度的分隔符将字符串分割成若干小块，
// 并返回一个 slice，如果字符串只包含空白符号，则返回一个长度为 0 的 slice。
//strings.Split(s, sep) 用于自定义分割符号来对指定字符串进行分割，同样返回 slice。
//因为这 2 个函数都会返回 slice，所以习惯使用 for-range 循环来对其进行处理（第 7.3 节）。
package main

import (
	"strings"
	"fmt"
)

func main() {
	var oriStr string = "abcdef ghijklmn opqrst uvwxyz"

	slice1 := strings.Fields(oriStr)
	fmt.Printf("slice1: %v\n", slice1) //slice1: [abcdef ghijklmn opqrst uvwxyz]

	oriStr = "abcdef-ghijklmn-opqrst-uvwxyz"

	slice2 := strings.Split(oriStr, "-")
	fmt.Printf("slice2: %v\n", slice2)
	for i, x := range slice2 {
		fmt.Printf("index %v value is %v\n", i, x)
	}

}
