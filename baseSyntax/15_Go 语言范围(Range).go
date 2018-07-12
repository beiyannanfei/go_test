//Go 语言中 range 关键字用于for循环中迭代数组(array)、切片(slice)、通道(channel)或集合(map)的元素。
// 在数组和切片中它返回元素的索引值，在集合中返回 key-value 对的 key 值
package main

import "fmt"

func main() {
	nums := []int{10, 20, 30, 40}
	sum := 0
	for key, value := range nums {
		fmt.Printf("key: %d, value: %d, nums[%d]=%d\n", key, value, key, nums[key])
		sum += value
	}
	fmt.Println("sum=", sum)
	/*
	key: 0, value: 10, nums[0]=10
	key: 1, value: 20, nums[1]=20
	key: 2, value: 30, nums[2]=30
	key: 3, value: 40, nums[3]=40
	sum= 100
	*/
	fmt.Println("---------------------------")
	//range也可以用在map的键值对上。
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}
	/*
	a -> apple
	b -> banana
	*/

	fmt.Println("---------------------------")
	//range也可以用来枚举Unicode字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
	for index, c := range "abcdef" {
		fmt.Printf("%d code is %d\n", index, c)
	}
	/*
	0 code is 97
	1 code is 98
	2 code is 99
	3 code is 100
	4 code is 101
	5 code is 102
	*/
}
