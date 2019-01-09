package main

import (
	"log"
	"bytes"
	"reflect"
)

//https://studygolang.com/articles/17367

func main() {
	makeSlice()
	appendSlice1()
	appendSlice2()
	appendSlice3()
	appendSlice4()
	rangeSlice()
	equalSlice()
}

func equalSlice() {
	log.Printf("===============equalSlice begin================\n\n")
	//equalSlice 数据切片的比较

	var slice1 []byte
	slice2 := []byte{}
	log.Println("out: slice1 == slice1: ", bytes.Equal(slice1, slice2))
	//out: slice1 == slice1:  true
	log.Println("out: slice1 == slice2: ", reflect.DeepEqual(slice1, slice2))
	//out: slice1 == slice2:  false

	log.Printf("===============equalSlice end================\n\n")
}

func rangeSlice() {
	log.Printf("===============rangeSlice begin================\n\n")

	slice := []int{10, 20, 30, 40,} //这个","不是必须的，但是如果换行之后这个","就是必须的

	//迭代每个元素,并显示值和地址，这个值是原来元素值的一份拷贝，修改这个值并不会改变原来元素的值
	for index, value := range slice {
		log.Printf("out: Value: %d Value-Addr: %X ElemAddr: %X\n", value, &value, &slice[index])
		//out: Value: 10 Value-Addr: C420070E58 ElemAddr: C420094060
		//out: Value: 20 Value-Addr: C420070E58 ElemAddr: C420094068
		//out: Value: 30 Value-Addr: C420070E58 ElemAddr: C420094070
		//out: Value: 40 Value-Addr: C420070E58 ElemAddr: C420094078
		value *= 10
	}
	log.Printf("out: slice = %v\n", slice)
	//out: slice = [10 20 30 40]

	//需要使用数据切片的引用才能实现对源数据的修改
	for index, value := range slice {
		log.Printf("out: Value: %d Value-Addr: %X ElemAddr: %X\n", value, &value, &slice[index])
		//out: Value: 10 Value-Addr: C420014F18 ElemAddr: C4200160C0
		//out: Value: 20 Value-Addr: C420014F18 ElemAddr: C4200160C8
		//out: Value: 30 Value-Addr: C420014F18 ElemAddr: C4200160D0
		//out: Value: 40 Value-Addr: C420014F18 ElemAddr: C4200160D8
		slice[index] *= 10
	}
	log.Printf("out: slice = %v\n", slice)
	//out: slice = [100 200 300 400]

	log.Printf("===============rangeSlice end================\n\n")
}

func appendSlice4() {
	log.Printf("===============appendSlice4 begin================\n\n")

	slice := [][]int{{10}, {100, 200}}
	log.Printf("out: slice = %v\t len(slice) == %d,\t cap(slice) == %d\n\n", slice, len(slice), cap(slice))
	//out: slice = [[10] [100 200]]	 len(slice) == 2,	 cap(slice) == 2

	log.Printf("out: slice[0] = %v\t len(slice[0]) == %d,\t cap(slice[0]) == %d\n\n", slice[0], len(slice[0]), cap(slice[0]))
	//out: slice[0] = [10]	 len(slice[0]) == 1,	 cap(slice[0]) == 1

	log.Printf("out: slice[1] = %v\t len(slice[1]) == %d,\t cap(slice[1]) == %d\n\n", slice[1], len(slice[1]), cap(slice[1]))
	//out: slice[1] = [100 200]	 len(slice[1]) == 2,	 cap(slice[1]) == 2

	slice[0] = append(slice[0], 20)
	log.Printf("out: slice = %v\t len(slice) == %d,\t cap(slice) == %d\n\n", slice, len(slice), cap(slice))
	//out: slice = [[10 20] [100 200]]	 len(slice) == 2,	 cap(slice) == 2

	log.Printf("out: slice[0] = %v\t len(slice[0]) == %d,\t cap(slice[0]) == %d\n\n", slice[0], len(slice[0]), cap(slice[0]))
	//out: slice[0] = [10 20]	 len(slice[0]) == 2,	 cap(slice[0]) == 2

	log.Printf("out: slice[1] = %v\t len(slice[1]) == %d,\t cap(slice[1]) == %d\n\n", slice[1], len(slice[1]), cap(slice[1]))
	//out: slice[1] = [100 200]	 len(slice[1]) == 2,	 cap(slice[1]) == 2

	log.Printf("===============appendSlice4 end================\n\n")
}

func appendSlice3() {
	log.Printf("===============appendSlice3 begin================\n\n")
	//数据切片声明的时候，直接可以通过append添加元素
	//验证数据切片的容量是按照倍数增长的

	var slice []int
	i := 0
	for i < 10 {
		log.Printf("out: slice = %v\t len(slice) == %d,\t cap(slice) == %d\n", slice, len(slice), cap(slice))
		//out: slice = []	 len(slice) == 0,	 cap(slice) == 0
		//out: slice = [0]	 len(slice) == 1,	 cap(slice) == 1
		//out: slice = [0 1]	 len(slice) == 2,	 cap(slice) == 2
		//out: slice = [0 1 2]	 len(slice) == 3,	 cap(slice) == 4
		//out: slice = [0 1 2 3]	 len(slice) == 4,	 cap(slice) == 4
		//out: slice = [0 1 2 3 4]	 len(slice) == 5,	 cap(slice) == 8
		//out: slice = [0 1 2 3 4 5]	 len(slice) == 6,	 cap(slice) == 8
		//out: slice = [0 1 2 3 4 5 6]	 len(slice) == 7,	 cap(slice) == 8
		//out: slice = [0 1 2 3 4 5 6 7]	 len(slice) == 8,	 cap(slice) == 8
		//out: slice = [0 1 2 3 4 5 6 7 8]	 len(slice) == 9,	 cap(slice) == 16
		slice = append(slice, i)
		i++
	}
	log.Printf("\n\t\t\t\t    out: slice = %v\t len(slice) == %d,\t cap(slice) == %d\n\n", slice, len(slice), cap(slice))
	//out: slice = [0 1 2 3 4 5 6 7 8 9]	 len(slice) == 10,	 cap(slice) == 16

	//直接使用 slice：即使函数内部得到的是 slice 的值拷贝，但依旧会更新 slice 的原始数据（底层 array）
	//会修改 slice 的底层 array，从而修改 slice
	func(arr []int) {
		arr[0] = 7
		log.Printf("out: slice = %v\n\n", slice)
		//out: slice = [7 1 2 3 4 5 6 7 8 9]
	}(slice)

	log.Printf("out: slice = %v\n\n", slice)
	//out: slice = [7 1 2 3 4 5 6 7 8 9]

	log.Printf("===============appendSlice3 end================\n\n")
}

func appendSlice2() {
	log.Printf("===============appendSlice2 begin================\n\n")
	//数据切片是一个引用，但是容量扩充之后引用就会发生变化

	parent_arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	log.Printf("out: parent_arr = %v\n\n", parent_arr)
	//out: parent_arr = [1 2 3 4 5 6 7 8 9 10]

	child1_slice := parent_arr[0:10]
	log.Printf("out: child1_slice = %v len(child1_slice) == %d,\t cap(child1_slice) == %d\n\n", child1_slice, len(child1_slice), cap(child1_slice))
	//out: child1_slice = [1 2 3 4 5 6 7 8 9 10] len(child1_slice) == 10,	 cap(child1_slice) == 10

	child2_slice := child1_slice[3:5]
	log.Printf("out: child2_slice = %v len(child2_slice) == %d,\t cap(child2_slice) == %d\n\n", child2_slice, len(child2_slice), cap(child2_slice))
	//out: child2_slice = [4 5] len(child2_slice) == 2,	 cap(child2_slice) == 7

	child2_slice = append(child2_slice, 100)
	log.Printf("out: child2_slice = %v len(child2_slice) == %d,\t cap(child2_slice) == %d\n\n", child2_slice, len(child2_slice), cap(child2_slice))
	//out: child2_slice = [4 5 100] len(child2_slice) == 3,	 cap(child2_slice) == 7

	log.Printf("out: child1_slice = %v len(child1_slice) == %d,\t cap(child1_slice) == %d\n\n", child1_slice, len(child1_slice), cap(child1_slice))
	//out: child1_slice = [1 2 3 4 5 100 7 8 9 10] len(child1_slice) == 10,	 cap(child1_slice) == 10

	child3_slice := parent_arr[3:5]
	log.Printf("out: child3_slice = %v len(child3_slice) == %d,\t cap(child3_slice) == %d\n\n", child3_slice, len(child3_slice), cap(child3_slice))
	//out: child3_slice = [4 5] len(child3_slice) == 2,	 cap(child3_slice) == 7

	child3_slice = append(child3_slice, 2000)
	log.Printf("out: child3_slice = %v len(child3_slice) == %d,\t cap(child3_slice) == %d\n\n", child3_slice, len(child3_slice), cap(child3_slice))
	//out: child3_slice = [4 5 2000] len(child3_slice) == 3,	 cap(child3_slice) == 7

	log.Printf("out: child1_slice = %v len(child1_slice) == %d,\t cap(child1_slice) == %d\n\n", child1_slice, len(child1_slice), cap(child1_slice))
	//out: child1_slice = [1 2 3 4 5 2000 7 8 9 10] len(child1_slice) == 10,	 cap(child1_slice) == 10

	log.Printf("out: child2_slice = %v len(child2_slice) == %d,\t cap(child2_slice) == %d\n\n", child2_slice, len(child2_slice), cap(child2_slice))
	//out: child2_slice = [4 5 2000] len(child2_slice) == 3,	 cap(child2_slice) == 7

	log.Printf("out: &child3_slice[1] = %v\t &child2_slice[1] = %v\n\n", &child3_slice[1], &child2_slice[1])
	//out: &child3_slice[1] = 0xc42001c3e0	 &child2_slice[1] = 0xc42001c3e0

	child4_slice := parent_arr[3:5:5] //s = s[low : high : max] 切片的三个参数的切片截取的意义为 low为截取的起始下标（含）， high为窃取的结束下标（不含high），max为切片保留的原切片的最大下标（不含max）；即新切片从老切片的low下标元素开始，len = high - low, cap = max - low；high 和 max一旦超出在老切片中越界，就会发生runtime err，slice out of range。另外如果省略第三个参数的时候，第三个参数默认和第二个参数相同，即len = cap
	log.Printf("out: child4_slice = %v len(child4_slice) == %d,\t cap(child4_slice) == %d\n\n", child4_slice, len(child4_slice), cap(child4_slice))
	//out: child4_slice = [4 5] len(child4_slice) == 2,	 cap(child4_slice) == 2

	child4_slice = append(child4_slice, 30000)
	log.Printf("out: child4_slice = %v len(child4_slice) == %d,\t cap(child4_slice) == %d\n\n", child4_slice, len(child4_slice), cap(child4_slice))
	//out: child4_slice = [4 5 30000] len(child4_slice) == 3,	 cap(child4_slice) == 4

	log.Printf("out: child1_slice = %v len(child1_slice) == %d,\t cap(child1_slice) == %d\n\n", child1_slice, len(child1_slice), cap(child1_slice))
	//out: child1_slice = [1 2 3 4 5 2000 7 8 9 10] len(child1_slice) == 10,	 cap(child1_slice) == 10

	log.Printf("out: child2_slice = %v len(child2_slice) == %d,\t cap(child2_slice) == %d\n\n", child2_slice, len(child2_slice), cap(child2_slice))
	//out: child2_slice = [4 5 2000] len(child2_slice) == 3,	 cap(child2_slice) == 7

	log.Printf("out: child3_slice = %v len(child3_slice) == %d,\t cap(child3_slice) == %d\n\n", child3_slice, len(child3_slice), cap(child3_slice))
	//out: child3_slice = [4 5 2000] len(child3_slice) == 3,	 cap(child3_slice) == 7

	child5_slice := parent_arr[3:5:6]
	log.Printf("out: child5_slice = %v len(child5_slice) == %d,\t cap(child5_slice) == %d\n\n", child5_slice, len(child5_slice), cap(child5_slice))
	//out: child5_slice = [4 5] len(child5_slice) == 2,	 cap(child5_slice) == 3

	child5_slice = append(child5_slice, 400000)
	log.Printf("out: child5_slice = %v len(child5_slice) == %d,\t cap(child5_slice) == %d\n\n", child5_slice, len(child5_slice), cap(child5_slice))
	//out: child5_slice = [4 5 400000] len(child5_slice) == 3,	 cap(child5_slice) == 3

	log.Printf("out: child1_slice = %v len(child1_slice) == %d,\t cap(child1_slice) == %d\n\n", child1_slice, len(child1_slice), cap(child1_slice))
	//out: child1_slice = [1 2 3 4 5 400000 7 8 9 10] len(child1_slice) == 10,	 cap(child1_slice) == 10

	log.Printf("out: child2_slice = %v len(child2_slice) == %d,\t cap(child2_slice) == %d\n\n", child2_slice, len(child2_slice), cap(child2_slice))
	//out: child2_slice = [4 5 400000] len(child2_slice) == 3,	 cap(child2_slice) == 7

	log.Printf("out: child3_slice = %v len(child3_slice) == %d,\t cap(child3_slice) == %d\n\n", child3_slice, len(child3_slice), cap(child3_slice))
	//out: child3_slice = [4 5 400000] len(child3_slice) == 3,	 cap(child3_slice) == 7

	log.Printf("out: child4_slice = %v len(child4_slice) == %d,\t cap(child4_slice) == %d\n\n", child4_slice, len(child4_slice), cap(child4_slice))
	//out: child4_slice = [4 5 30000] len(child4_slice) == 3,	 cap(child4_slice) == 4

	child5_slice = append(child5_slice, 5000000)
	log.Printf("out: child5_slice = %v len(child5_slice) == %d,\t cap(child5_slice) == %d\n\n", child5_slice, len(child5_slice), cap(child5_slice))
	//out: child5_slice = [4 5 400000 5000000] len(child5_slice) == 4,	 cap(child5_slice) == 6

	log.Printf("out: child1_slice = %v len(child1_slice) == %d,\t cap(child1_slice) == %d\n\n", child1_slice, len(child1_slice), cap(child1_slice))
	//out: child1_slice = [1 2 3 4 5 400000 7 8 9 10] len(child1_slice) == 10,	 cap(child1_slice) == 10

	log.Printf("out: child2_slice = %v len(child2_slice) == %d,\t cap(child2_slice) == %d\n\n", child2_slice, len(child2_slice), cap(child2_slice))
	//out: child2_slice = [4 5 400000] len(child2_slice) == 3,	 cap(child2_slice) == 7

	log.Printf("out: child3_slice = %v len(child3_slice) == %d,\t cap(child3_slice) == %d\n\n", child3_slice, len(child3_slice), cap(child3_slice))
	//out: child3_slice = [4 5 400000] len(child3_slice) == 3,	 cap(child3_slice) == 7

	log.Printf("out: child4_slice = %v len(child4_slice) == %d,\t cap(child4_slice) == %d\n\n", child4_slice, len(child4_slice), cap(child4_slice))
	//out: child4_slice = [4 5 30000] len(child4_slice) == 3,	 cap(child4_slice) == 4

	log.Printf("===============appendSlice2 end================\n\n")
}

func appendSlice1() {
	log.Printf("===============appendSlice1 begin================\n\n")

	oldSlice := make([]int, 10, 100)
	log.Printf("out: oldSlice = %v\t len(oldSlice) == %d\t cap(oldSlice) == %d\n\n", oldSlice, len(oldSlice), cap(oldSlice))
	//out: oldSlice = [0 0 0 0 0 0 0 0 0 0]	 len(oldSlice) == 10	 cap(oldSlice) == 100

	n, m := 1, 99
	newSlice := oldSlice[n:m]
	log.Printf("out: newSlice = %v\t len(newSlice) == %d\t cap(newSlice) == %d\n\n", newSlice, len(newSlice), cap(newSlice))
	//out: newSlice = [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]	 len(newSlice) == 98	 cap(newSlice) == 99

	log.Printf("out: len(newSlice) == m-n == %d, cap(newSlice) == cap(oldSlice)-n == %d\n\n", len(newSlice), cap(newSlice))
	//out: len(newSlice) == m-n == 98, cap(newSlice) == cap(oldSlice)-n == 99

	log.Printf("===============appendSlice1 end================\n\n")
}

func makeSlice() {
	log.Printf("===============makeSlice begin================\n\n")

	//声明数据切片 可以直接使用，直接使用append方法
	var slice1 []int
	log.Printf("out: slice1 = %v len(slice1) == %d,\t cap(slice1) == %d\n\n", slice1, len(slice1), cap(slice1))
	//out: slice1 = [] len(slice1) == 0,	 cap(slice1) == 0

	//初始化数据切片，指定数据切片的长度
	slice2 := make([]int, 10)
	log.Printf("out: slice2 = %v len(slice2) == %d,\t cap(slice2) == %d\n\n", slice2, len(slice2), cap(slice2))
	//out: slice2 = [0 0 0 0 0 0 0 0 0 0] len(slice2) == 10,	 cap(slice2) == 10

	//初始化数据切片，指定数据切片的长度和切片的最大容量
	slice3 := make([]int, 10, 100)
	log.Printf("out: slice3 = %v len(slice3) == %d,\t cap(slice3) == %d\n\n", slice3, len(slice3), cap(slice3))
	//out: slice3 = [0 0 0 0 0 0 0 0 0 0] len(slice3) == 10,	 cap(slice3) == 100

	//声明数据切片并赋值数据切片，此时数据切片的长度和最大容量相等
	slice4 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	log.Printf("out: slice4 = %v len(slice4) == %d,\t cap(slice4) == %d\n\n", slice4, len(slice4), cap(slice4))
	//out: slice4 = [1 2 3 4 5 6 7 8 9 0] len(slice4) == 10,	 cap(slice4) == 10

	slice5 := append(slice4, 10)
	log.Printf("out: slice5 = %v len(slice5) == %d,\t cap(slice5) == %d\n\n", slice5, len(slice5), cap(slice5))
	//out: slice5 = [1 2 3 4 5 6 7 8 9 0 10] len(slice5) == 11,	 cap(slice5) == 20

	log.Printf("===============makeSlice end================\n\n")
}
