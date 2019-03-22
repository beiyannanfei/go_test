package main

import "fmt"

// 13~23 例子来源 https://juejin.im/post/5b5bd2ddf265da0f716c2fea?utm_source=gold_browser_extension

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
