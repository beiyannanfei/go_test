//Go 语言切片是对数组的抽象。
//Go 数组的长度不可改变，在特定场景中这样的集合就不太适用，Go中提供了一种灵活，功能强悍的内置类型切片("动态数组"),与数组相比切片的长度是不固定的，可以追加元素，在追加时可能使切片的容量增大。

package main

import "fmt"

func main() {
	//s := [] int{1, 2, 3}	//直接初始化切片，[]表示是切片类型，{1,2,3}初始化值依次是1,2,3.其cap=len=3
	//s :=make([]int,len,cap)	//通过内置函数make()初始化切片s,[]int 标识为其元素类型为int的切片
	printSlice()
	nilSlice()
	fmt.Println("-------------------------------")
	subSlice()
	fmt.Println("===============================")
	appendSlice()
	fmt.Println("-------------------------------")
	copySlice()
}

func copySlice() {
	number := make([]int, 0, 3)
	fmt.Println(number, cap(number)) //[] 3
	number = append(number, 0, 1, 2, 3)
	fmt.Println(number, cap(number))                   //[0 1 2 3] 6
	number1 := make([]int, len(number), cap(number)*2) //创建切片 numbers1 是之前切片的两倍容量
	fmt.Println(number1, len(number1), cap(number1))   //[0 0 0 0] 4 12
	copy(number1, number)                              //拷贝 number 的内容到 numbers1
	fmt.Println(number1, len(number1), cap(number1))   //[0 1 2 3] 4 12
}

func appendSlice() { //切片追加元素
	var numbers [] int
	fmt.Println("原始切片numbers =", numbers)                                               //原始切片numbers = []
	numbers = append(numbers, 0)                                                        //向空切片追加元素
	fmt.Println("numbers =", numbers)                                                   //numbers = [0]
	numbers = append(numbers, 1)                                                        //添加一个元素
	fmt.Println("numbers =", numbers)                                                   //numbers = [0 1]
	numbers = append(numbers, 2, 3, 4)                                                  //同时追加多个元素
	fmt.Println("numbers =", numbers)                                                   //numbers = [0 1 2 3 4]
	fmt.Printf("len = %d, cat = %d, numbers = %x\n", len(numbers), cap(numbers), numbers) //len = 5, cat = 6, numbers = [0 1 2 3 4]
}

func subSlice() { //切片截取
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	//len = 9 cap = 9 slice = [0 1 2 3 4 5 6 7 8]
	fmt.Printf("len = %d cap = %d slice = %v\n", len(numbers), cap(numbers), numbers)
	//原始切片 ==  [0 1 2 3 4 5 6 7 8]
	fmt.Println("原始切片 == ", numbers)
	//切片索引从1到4(不包含)数据 ==  [1 2 3]
	fmt.Println("切片索引从1到4(不包含)数据 == ", numbers[1:4])
	//切片索引3(不包含)数据 ==  [0 1 2] (默认下限为0)
	fmt.Println("切片到索引3(不包含)数据 == ", numbers[:3])
	//切片从索引4数据 ==  [4 5 6 7 8] (默认上限为len(slice))
	fmt.Println("切片从索引4数据 == ", numbers[4:])
	fmt.Println("-------------------------------")
	numbers1 := numbers[:2]
	fmt.Printf("len = %d cap = %d slice = %v\n", len(numbers1), cap(numbers1), numbers1) //len = 2 cap = 9 slice = [0 1]
	number2 := numbers[2:5]
	fmt.Printf("len = %d cap = %d slice = %v\n", len(number2), cap(number2), number2) //len = 3 cap = 7 slice = [2 3 4]
}

func nilSlice() { //空切片
	var numbers []int
	if (numbers == nil) {
		fmt.Printf("切片是空的\n")
	}
	fmt.Printf("len = %d cap = %d slice = %v\n", len(numbers), cap(numbers), numbers)
}

func printSlice() {
	var numbers = make([]int, 3, 5)
	//长度	容量		具体值
	fmt.Printf("len = %d cap = %d slice = %v\n", len(numbers), cap(numbers), numbers) //len = 3 cap = 5 slice = [0 0 0]
}
