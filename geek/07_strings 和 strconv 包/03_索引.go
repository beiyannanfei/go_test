//Index 返回字符串 str 在字符串 s 中的索引（str 的第一个字符的索引），-1 表示字符串 s 不包含字符串 str：
//strings.Index(s, str string) int

//LastIndex 返回字符串 str 在字符串 s 中最后出现位置的索引（str 的第一个字符的索引），-1 表示字符串 s 不包含字符串 str：
//strings.LastIndex(s, str string) int

//如果 ch 是非 ASCII 编码的字符，建议使用IndexRune函数来对字符进行定位
//strings.IndexRune(s string, r rune) int

package main

import (
	"fmt"
	"strings"
)

func main() {
	var str string = "Hi, I'm Marc, Hi."
	fmt.Printf("the position of 'Marc' is: %d\n", strings.Index(str, "Marc"))          //the position of 'Marc' is: 8
	fmt.Printf("the position of the last 'Hi' is: %d\n", strings.LastIndex(str, "Hi")) //the position of the last 'Hi' is: 14

	fmt.Printf("the position of 'abcdef' is: %d\n", strings.Index(str, "abcdef")) //the position of 'abcdef' is: -1

	var chStr string = "侯卫东官场笔记第二部有声小说官场"
	fmt.Printf("the position of '笔记' is: %d\n", strings.IndexRune(chStr, '笔'))       //the position of '笔记' is: 15
	fmt.Printf("the position of '笔记' is: %d\n", strings.IndexRune(chStr, rune('笔'))) //the position of '笔记' is: 15
}
