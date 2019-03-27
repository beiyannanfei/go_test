package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	params := make(map[string]string)

	params["id"] = "1"
	params["id1"] = "3"
	params["controller"] = "sections"
	params["Controller"] = "Sections"

	data, _ := json.Marshal(params)
	fmt.Println(string(data)) //{"Controller":"Sections","controller":"sections","id":"1","id1":"3"}
	//todo 利用Golang自带的json转换包转换，会将map中key的顺序改为字母顺序，而不是map的赋值顺序。map这个结构哪怕利用for range遍历的时候,其中的key也是无序的，可以理解为map就是个无序的结构，和php中的array要区分开来
}
