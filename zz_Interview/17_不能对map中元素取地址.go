package main

import "fmt"

// 13~23 例子来源 https://juejin.im/post/5b5bd2ddf265da0f716c2fea?utm_source=gold_browser_extension

func main() {
	ages := make(map[string]int)
	ages["aaa"] = 100

	fmt.Println(ages)
	//a := &ages["aaa"]	//cannot take the address of ages["aaa"]
	//todo map中的元素不是一个变量，不能对map的元素进行取地址操作，禁止对map进行取地址操作的原因可能是map随着元素的增加map可能会重新分配内存空间，这样会导致原来的地址无效
}
