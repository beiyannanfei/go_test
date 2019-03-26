package main

import "fmt"

func arr1() {
	array := [4]int{10, 20, 30, 40}
	slice := array[0:2]
	newSlice := append(slice, 50)
	newSlice[1] += 1
	fmt.Println(newSlice) //[10 21 50]
	fmt.Println(slice)    //[10 21]
}

func arr2() {
	array := [4]int{10, 20, 30, 40}
	slice := array[0:2]
	newSlice := append(append(append(slice, 50), 100), 150)
	newSlice[1] += 1
	fmt.Println(newSlice) //[10 21 50 100 150]
	fmt.Println(slice)    //[10 20]
}

func main() {
	arr1()
	fmt.Println("------------------------")
	arr2()
}
//todo 这就要从Golang切片的扩容说起了；切片的扩容，就是当切片添加元素时，切片容量不够了，就会扩容，扩容的大小遵循下面的原则：（如果切片的容量小于1024个元素，那么扩容的时候slice的cap就翻番，乘以2；一旦元素个数超过1024个元素，增长因子就变成1.25，即每次增加原来容量的四分之一。）如果扩容之后，还没有触及原数组的容量，那么，切片中的指针指向的位置，就还是原数组（这就是产生bug的原因）；如果扩容之后，超过了原数组的容量，那么，Go就会开辟一块新的内存，把原来的值拷贝过来，这种情况丝毫不会影响到原数组
