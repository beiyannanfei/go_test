package main

import "fmt"

func main() {
	var a int = 20 /* 声明实际变量 */
	var ip *int    /* 声明指针变量 */

	ip = &a /* 指针变量的存储地址 */

	fmt.Printf("a 变量的地址是: %x\n", &a)

	/* 指针变量的存储地址 */
	fmt.Printf("ip 变量储存的指针地址: %x\n", ip)

	/* 使用指针访问值 */
	fmt.Printf("*ip 变量的值: %d\n", *ip)

	nilPointer()
	fmt.Println("-----------------------")
	pointerArr()
}

func nilPointer() { //空指针
	var ptr *int
	fmt.Printf("ptr value is: %x\n", ptr)
	if (ptr != nil) {
		fmt.Println("非空指针")
	}
}

func pointerArr() { //指针数组
	arr := []int{10, 20, 30}
	var ptr [3]*int
	for i := 0; i < len(arr); i++ {
		ptr[i] = &arr[i]
	}
	for i := 0; i < len(ptr); i++ {
		fmt.Printf("arr[%d] = %d\n", i, *ptr[i])
	}
}





