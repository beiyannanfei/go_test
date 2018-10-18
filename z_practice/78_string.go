package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=================== func Contains(s, substr string) bool 字符串s中是否包含substr，返回bool值 ===================")
	fmt.Println(strings.Contains("seafood", "foo")) //=> true
	fmt.Println(strings.Contains("seafood", "bar")) //=>false

	fmt.Println("=================== func Join(a []string, sep string) string 字符串链接，把slice a通过sep链接起来 ===================")
	s := []string{"foo", "bar", "baz"}
	fmt.Println(strings.Join(s, ", ")) //=>foo, bar, baz

	fmt.Println("=================== func Index(s, sep string) int 在字符串s中查找sep所在的位置，返回位置值，找不到返回-1 ===================")
	fmt.Println(strings.Index("chicken", "ken")) //=>4
	fmt.Println(strings.Index("chicken", "dmr")) //=>-1

	fmt.Println("=================== func Repeat(s string, count int) string 重复s字符串count次，最后返回重复的字符串 ===================")
	fmt.Println("ba" + strings.Repeat("na", 2)) //=>banana

	fmt.Println("=================== func Replace(s, old, new string, n int) string 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换 ===================")
	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))      //=>oinky oinky oink
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1)) //=>moo moo moo

	fmt.Println("=================== func Split(s, sep string) []string 把s字符串按照sep分割，返回slice ===================")
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))                        //=>["a" "b" "c"]
	fmt.Printf("%q\n", strings.Split("a man a plan a canal panama", "a ")) //=>["" "man " "plan " "canal panama"]
	fmt.Printf("%q\n", strings.Split(" xyz ", ""))                         //=>[" " "x" "y" "z" " "]
	fmt.Printf("%q\n", strings.Split("", "a"))                             //=>[""]

	fmt.Println("=================== func Trim(s string, cutset string) string 在s字符串中去除cutset指定的字符串 ===================")
	fmt.Printf("[%q]\n", strings.Trim(" !!! Achtung  !!! ", " !"))  //=>["Achtung"]
	fmt.Printf("[%q]\n", strings.Trim("ababbaAchtungababba", "ab")) //=>["Achtung"]

	fmt.Println("=================== func Fields(s string) []string 去除s字符串的空格符，并且按照空格分割返回slice ===================")
	fmt.Printf("Fields are: %q", strings.Fields("   foo bar   baz     ")) //=>Fields are: ["foo" "bar" "baz"]
}
