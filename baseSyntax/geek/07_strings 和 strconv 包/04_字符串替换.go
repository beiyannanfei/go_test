//Replace 用于将字符串 str 中的前 n 个字符串 old 替换为字符串 new，并返回一个新的字符串，如果 n = -1 则替换所有字符串 old 为字符串 new：
//strings.Replace(str, old, new, n) string

package main

import (
	"strings"
	"fmt"
)

func main() {
	var str string = "aa bcd aa edf aa rfv aa gth aa tyfg aa"
	var srcStr string = "aa"
	var destStr string = "AA"
	var newStr = strings.Replace(str, srcStr, destStr, -1) //全部替换(注意：函数不会改变原字符串)
	fmt.Println(str)                                       //aa bcd aa edf aa rfv aa gth aa tyfg aa
	fmt.Println(newStr)                                    //AA bcd AA edf AA rfv AA gth AA tyfg AA

	newStr = strings.Replace(str, srcStr, destStr, 2) //只替换前2个
	fmt.Println(newStr)                               //AA bcd AA edf aa rfv aa gth aa tyfg aa
}
