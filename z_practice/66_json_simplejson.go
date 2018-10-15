package main

import (
	"github.com/bitly/go-simplejson"
	"fmt"
)

func main() {
	json, err := simplejson.NewJson([]byte(`{
		"test": {
			"array": [1, "2", 3],
			"int": 10,
			"float": 5.150,
			"bignum": 9223372036854775807, 
			"string": "simplejson", 
			"bool": true
		}
	}`))
	if err != nil {
		fmt.Println("simplejson.NewJson err:", err)
		return
	}

	fmt.Println("json:", json)

	arr, _ := json.Get("test").Get("array").Array()
	fmt.Println("arr:", arr)

	i, _ := json.Get("test").Get("int").Int()
	fmt.Println("i:", i)

	f, _ := json.Get("test").Get("float").Float64()
	fmt.Println("f:", f)

	i64, _ := json.Get("test").Get("bignum").Int64()
	fmt.Println("i64:", i64)

	s, _ := json.Get("test").Get("string").String()
	fmt.Println("s:", s)

	b, _ := json.Get("test").Get("bool").Bool()
	fmt.Println("b:", b)

	_, err = json.Get("test").Get("aaa").String()
	if err != nil {
		fmt.Println("json get err:", err)
	}
}
