package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("============== Append 系列函数将整数等转换为字符串后，添加到现有的字节数组中 ==============")
	str := make([]byte, 0, 100)
	str = strconv.AppendInt(str, 4567, 10)
	str = strconv.AppendBool(str, false)
	str = strconv.AppendQuote(str, "abcdefg")
	str = strconv.AppendQuoteRune(str, '单')
	fmt.Printf("%v\n", string(str))

	fmt.Println("============== Format 系列函数把其他类型的转换为字符串 ==============")
	a := strconv.FormatBool(false)
	b := strconv.FormatFloat(123.12, 'g', 12, 64)
	c := strconv.FormatInt(1234, 10)
	d := strconv.FormatUint(12345, 10)
	e := strconv.Itoa(1023)
	fmt.Printf("%#v %#v %#v %#v %#v\n", a, b, c, d, e)

	fmt.Println("============== Parse 系列函数把字符串转换为其他类型 ==============")
	aa, err := strconv.ParseBool("false")
	if err != nil {
		fmt.Println("strconv.ParseBool err", err)
	}

	bb, err := strconv.ParseFloat("123.23", 64)
	if err != nil {
		fmt.Println("strconv.ParseFloat err", err)
	}

	cc, err := strconv.ParseInt("1234", 10, 64)
	if err != nil {
		fmt.Println("strconv.ParseInt err", err)
	}

	dd, err := strconv.ParseUint("12345", 10, 64)
	if err != nil {
		fmt.Println("strconv.ParseUint err", err)
	}

	ee, err := strconv.Atoi("1023")
	fmt.Printf("%#v %#v %#v %#v %#v\n", aa, bb, cc, dd, ee)
}
