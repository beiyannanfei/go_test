//HasPrefix 判断字符串 s 是否以 prefix 开头：
//strings.HasPrefix(s, prefix string) bool

//HasSuffix 判断字符串 s 是否以 suffix 结尾：
//strings.HasSuffix(s, suffix string) bool
package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "This is an example of a string"
	fmt.Printf("变量str是否以Th开头？%v\n", strings.HasPrefix(str, "Th")) //变量str是否以Th开头？true
	fmt.Printf("变量str是否以ng结尾？%t", strings.HasSuffix(str, "ng"))   //变量str是否以ng结尾？true
}
