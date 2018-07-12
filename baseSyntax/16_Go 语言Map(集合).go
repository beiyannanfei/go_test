//Map 是一种无序的键值对的集合。Map 最重要的一点是通过 key 来快速检索数据，key 类似于索引，指向数据的值。
//Map 是一种集合，所以我们可以像迭代数组和切片那样迭代它。不过，Map 是无序的，我们无法决定它的返回顺序，这是因为 Map 是使用 hash 表来实现的。

/* 声明变量，默认 map 是 nil */
//var map_variable map[key_data_type]value_data_type

/* 使用 make 函数 */
//map_variable := make(map[key_data_type]value_data_type)

package main

import "fmt"

func main() {
	map1()
	fmt.Println("==============================")
	delMap()
}

func delMap() {
	countryCapitalMap := map[string]string{
		"France": "Paris",
		"Italy":  "Rome",
		"Japan":  "Tokyo",
		"India":  "New Delhi"}

	fmt.Println("原始 map")
	for country, capital := range countryCapitalMap {
		fmt.Printf("country: %s, capital: %s\n", country, capital)
	}
	//============all countryCapitalMap is: map[string]string{"India":"New Delhi", "France":"Paris", "Italy":"Rome", "Japan":"Tokyo"}
	fmt.Printf("============all countryCapitalMap is: %#v\n", countryCapitalMap)
	//删除元素
	delete(countryCapitalMap, "France")
	fmt.Println("删除元素后 map")

	for country := range countryCapitalMap {
		fmt.Printf("country: %s, capital: %s\n", country, countryCapitalMap[country])
	}
	//============all countryCapitalMap is: map[string]string{"Japan":"Tokyo", "India":"New Delhi", "Italy":"Rome"}
	fmt.Printf("============all countryCapitalMap is: %#v\n", countryCapitalMap)
}

func map1() {
	var countryCapitalMap map[string]string
	/* 创建集合 */
	countryCapitalMap = make(map[string]string)

	/* map 插入 key-value 对，各个国家对应的首都 */
	countryCapitalMap["France"] = "Paris"
	countryCapitalMap["Italy"] = "Rome"
	countryCapitalMap["Japan"] = "Tokyo"
	countryCapitalMap["India"] = "New Delhi"

	for index := range countryCapitalMap {
		fmt.Printf("country %s, capital is %s\n", index, countryCapitalMap[index])
	}
	fmt.Println("---------------------------")
	var capitial string
	var ok bool
	capitial, ok = countryCapitalMap["France"]
	fmt.Println(capitial, ok) //Paris true
	capitial, ok = countryCapitalMap["United States"]
	fmt.Println(capitial, ok) // false
	if ok {
		fmt.Println("Capital of United States is", capitial)
	} else {
		fmt.Println("Capital of United States is not present") //Capital of United States is not present
	}
}
