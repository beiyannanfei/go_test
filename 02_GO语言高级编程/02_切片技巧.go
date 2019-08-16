package main

import (
	"fmt"
)

func main() {
	defineSlice()
	appendSlice()
	unshiftSlice()
	addSliceWithCopy()
	deleteSlice()
}

func deleteSlice() {
	fmt.Println("================== 删除切片元素 ==================")
	fmt.Println("***删除尾部***")
	var a = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var index, n = 3, 3
	a = a[:len(a)-1]               //删除尾部1个元素
	fmt.Println(a, len(a), cap(a)) //[0 1 2 3 4 5 6 7 8] 9 10
	a = a[:len(a)-n]               //删除尾部n个元素
	fmt.Println(a, len(a), cap(a)) //[0 1 2 3 4 5] 6 10

	fmt.Println("***删除头部(切边位置方式实现)***")
	a = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = a[1:]                      //删除头部1个元素
	fmt.Println(a, len(a), cap(a)) //[1 2 3 4 5 6 7 8 9] 9 9
	a = a[n:]                      //删除头部n个元素
	fmt.Println(a, len(a), cap(a)) //[4 5 6 7 8 9] 6 6

	fmt.Println("***删除头部(append方式实现)***")
	a = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = append(a[:0], a[1:]...)    //删除开头1个元素
	fmt.Println(a, len(a), cap(a)) //[1 2 3 4 5 6 7 8 9] 9 10
	a = append(a[:0], a[n:]...)    //删除开头n个元素
	fmt.Println(a, len(a), cap(a)) //[4 5 6 7 8 9] 6 10

	fmt.Println("***删除头部(copy方式实现)***")
	a = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = a[:copy(a, a[1:])]         //删除开头1个元素
	fmt.Println(a, len(a), cap(a)) //[1 2 3 4 5 6 7 8 9] 9 10
	a = a[:copy(a, a[n:])]         //删除开头n个元素
	fmt.Println(a, len(a), cap(a)) //[4 5 6 7 8 9] 6 10

	fmt.Println("***删除中间(append方式实现)***")
	a = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	a = append(a[:index], a[index+1:]...) //删除中间1个元素
	fmt.Println(a, len(a), cap(a))        //[0 1 2 4 5 6 7 8 9] 9 10
	a = append(a[:index], a[index+n:]...) //删除中间n个元素
	fmt.Println(a, len(a), cap(a))        //[0 1 2 7 8 9] 6 10

}

func addSliceWithCopy() {
	fmt.Println("================== 利用copy函数实现切片插入元素 ==================")
	var a = []int{1, 2, 3, 4, 5}
	var i, x = 2, 10000
	a = append(a, 0)               //切片扩展1个空间
	copy(a[i+1:], a[i:])           //a[i:]向后移动1个位置
	a[i] = x                       //设置新添加的元素
	fmt.Println(a, len(a), cap(a)) //[1 2 10000 3 4 5] 6 10

	var b = []int{1, 2, 3, 4, 5}
	var temp = []int{10000, 20000, 30000}
	b = append(b, temp...)         //为切片扩展足够的空间
	copy(b[i+len(temp):], b[i:])   //a[i:]向后移动len(temp)个位置
	copy(b[i:], temp)              //复制新添加的切边
	fmt.Println(b, len(b), cap(b)) //[1 2 10000 20000 30000 3 4 5] 8 10
}

func unshiftSlice() {
	fmt.Println("================== 切片头部添加元素 ==================")
	var a = []int{1, 2, 3}
	a = append([]int{0}, a...)          //在开头添加1个元素
	fmt.Println(a, len(a), cap(a))      //[0 1 2 3] 4 4
	a = append([]int{-3, -2, -1}, a...) //在开头添加1个切片
	fmt.Println(a, len(a), cap(a))      //[-3 -2 -1 0 1 2 3] 7 8

	fmt.Println("================== 切片中间添加元素 ==================")
	var b = []int{1, 2, 3, 4, 5}
	var i, x = 2, 100000
	b = append(b[:i], append([]int{x}, b[i:]...)...) //在第i个位置插入x
	fmt.Println(b, len(b), cap(b))                   //[1 2 100000 3 4 5] 6 10

	var c = []int{10, 20, 30, 40, 50}
	c = append(c[:i], append([]int{1, 2, 3}, c[i:]...)...) //在第i个位置插入切片
	fmt.Println(c, len(c), cap(c))                         //[10 20 1 2 3 30 40 50] 8 10
}

func appendSlice() {
	fmt.Println("================== 添加切片元素 ==================")
	var a []int
	a = append(a, 1)                 //追加1个元素
	fmt.Println(a, len(a), cap(a))   //[1] 1 1
	a = append(a, 1, 2, 3)           //追加多个元素，手写解包方式
	fmt.Println(a, len(a), cap(a))   //[1 1 2 3] 4 4
	a = append(a, []int{1, 2, 3}...) //追加一个切片，切片需要解包
	fmt.Println(a, len(a), cap(a))   //[1 1 2 3 1 2 3] 7 8
}

func defineSlice() {
	fmt.Println("================== 切片定义 ==================")
	var a []int                           //nil切片，和nil相等，一般用来表示一个不存在的切片
	fmt.Println(a == nil, len(a), cap(a)) //true 0 0

	var b = []int{}                       //空切片，和nil不相等，一般用来表示一个空集合
	fmt.Println(b == nil, len(a), cap(a)) //false 0 0

	var c = []int{1, 2, 3}         //有3个元素的切片，len和cap都为3
	fmt.Println(c, len(c), cap(c)) //[1 2 3] 3 3

	var d = c[:2]                  //有2个元素的切片，len为2， cap为3
	fmt.Println(d, len(d), cap(d)) //[1 2] 2 3

	var e = c[0:2:cap(c)]          //有两个元素的切片，len为2，cap为3
	fmt.Println(e, len(e), cap(e)) //[1 2] 2 3

	var f = c[:0]                  //有0个元素的切片， len为0， cap为3
	fmt.Println(f, len(f), cap(f)) //[] 0 3

	var g = make([]int, 3)         //有3个元素的切片，len和cap都为3
	fmt.Println(g, len(g), cap(g)) //[0 0 0] 3 3

	var h = make([]int, 2, 3)      //有2个元素的切片，len为2， cap为3
	fmt.Println(h, len(h), cap(h)) //[0 0] 2 3

	var i = make([]int, 0, 3)      //有0个元素的切片，len为0，cap为3
	fmt.Println(i, len(i), cap(i)) //[] 0 3
}
