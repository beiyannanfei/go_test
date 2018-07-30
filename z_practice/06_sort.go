package main

import (
	"sort"
	"fmt"
)

func main() {
	intList := []int{2, 4, 3, 5, 7, 6, 9, 8, 1, 0}
	oriIntList := make([]int, len(intList), cap(intList))
	copy(oriIntList, intList)
	sort.Ints(intList) //升序		Float64s  Strings
	fmt.Printf("oriIntList: %#v, intList: %#v\n", oriIntList, intList)
	sort.Sort(sort.Reverse(sort.IntSlice(intList)))	//降序
	fmt.Printf("oriIntList: %#v, intList: %#v\n", oriIntList, intList)
}

