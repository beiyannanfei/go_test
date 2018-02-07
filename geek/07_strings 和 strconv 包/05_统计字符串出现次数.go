//Count 用于计算字符串 str 在字符串 s 中出现的非重叠次数：
//strings.Count(s, str string) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "Hello, how is it going, Hugo?"
	var manyG = "gggggggggg"

	fmt.Printf("Number of H's in %s is: ", str)
	fmt.Printf("%d\n", strings.Count(str, "H")) //Number of H's in Hello, how is it going, Hugo? is: 2

	fmt.Printf("Number of double g's in %s is: ", manyG)
	fmt.Printf("%d\n", strings.Count(manyG, "gg")) //Number of double g's in gggggggggg is: 5
}
