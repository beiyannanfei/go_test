//Contains 判断字符串 s 是否包含 substr：
//strings.Contains(s, substr string) bool

package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "This is an example of a string"
	fmt.Printf("字符串str是否包含子串is？%v\n", strings.Contains(str, "is")) //字符串str是否包含子串is？true
	fmt.Printf("字符串str是否包含子串ab？%v\n", strings.Contains(str, "ab")) //字符串str是否包含子串ab？false
}
