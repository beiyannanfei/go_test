//你可以使用 strings.TrimSpace(s) 来剔除字符串开头和结尾的空白符号；
// 如果你想要剔除指定字符，则可以使用 strings.Trim(s, "cut") 来将开头和结尾的 cut 去除掉。
// 该函数的第二个参数可以包含任何字符，如果你只想剔除开头或者结尾的字符串，
// 则可以使用 TrimLeft 或者 TrimRight 来实现。
package main

import (
	"fmt"
	"strings"
)

func main() {
	var oriStr1 = "   aaa bbb   "
	fmt.Printf("oriStr1 is : %v\n", oriStr1)		//oriStr1 is :    aaa bbb
	var trimStr1 = strings.TrimSpace(oriStr1)
	fmt.Printf("trimStr1 is : %v\n", trimStr1)	//trimStr1 is : aaa bbb

	var oriStr2 = "aaaBBBBCCCCaaaa"
	fmt.Printf("oriStr2 is : %v\n", oriStr2)		//oriStr2 is : aaaBBBBCCCCaaaa
	var trimStr2 = strings.Trim(oriStr2, "a")
	fmt.Printf("trimStr2 is : %v\n", trimStr2)	//trimStr2 is : BBBBCCCC

	var leftTrimStr2 = strings.TrimLeft(oriStr2, "a")
	fmt.Printf("leftTrimStr2 is : %v\n", leftTrimStr2)	//leftTrimStr2 is : BBBBCCCCaaaa

	var rightTrimStr2 = strings.TrimRight(oriStr2, "a")
	fmt.Printf("rightTrimStr2 is : %v\n", rightTrimStr2)	//rightTrimStr2 is : aaaBBBBCCCC
}
