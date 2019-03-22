package main

import "fmt"

func main() {

	var sampleMap map[string]int

	sampleMap["test"] = 1	//assignment to entry in nil map
	//todo 必须使用make或者将map初始化之后，才可以添加元素。
	fmt.Println(sampleMap)

	/*
		var sampleMap map[string]int
    	sampleMap = map[string]int { //todo 初始化
    	    "test1":1,
    	}
    	sampleMap["test"] = 1
    	fmt.Println(sampleMap)
	*/
}
