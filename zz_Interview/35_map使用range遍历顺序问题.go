package main

import "fmt"

func main() {
	newMap := make(map[int]int)
	for i := 0; i < 10; i++ {
		newMap[i] = i
	}

	for key, value := range newMap {
		fmt.Printf("key is %d, value is %d\n", key, value)
	}
	//todo map使用range遍历顺序问题，并不是录入的顺序，而是随机顺序 是杂乱无章的顺序。map的遍历顺序不固定，这种设计是有意为之的，能为能防止程序依赖特定遍历顺序。
}
