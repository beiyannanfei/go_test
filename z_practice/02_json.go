package main

import (
	"fmt"
	"reflect"
	"encoding/json"
)
// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 02_json.go
func main() {
	type Person struct {		//注意：首字母要大写
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Addr  string `json:"addr"`
		Score []int  `json:"-"`		//-号代表不解析这个字段
	}

	p1 := Person{
		Name:  "zhangsan",
		Age:   25,
		Addr:  "beijing",
		Score: []int{10, 20, 30},
	}

	p1Str, err := json.Marshal(p1) //类似JSON.stringify
	if err != nil {
		fmt.Printf("json.Marshal(p1) err: %v, p1: %v\n", err.Error(), p1)
		return
	}
	fmt.Printf("p1Str: %#v, p1Str type: %v, p1: %v\n", p1Str, reflect.TypeOf(p1Str).String(), p1)

	var p2 Person
	json.Unmarshal(p1Str, &p2) //类似JSON.parse
	fmt.Printf("p2: %v\n", p2)

	type Student struct {
		P      Person
		School string
	}

	s1 := Student{
		P: Person{
			Name:  "lisi",
			Age:   18,
			Addr:  "hebei",
			Score: []int{11, 22, 33},
		},
		School: "hbgydx",
	}
	s1Str, _ := json.Marshal(s1)
	fmt.Printf("s1Str: %#v, s1Str type: %v, s1: %v\n", s1Str, reflect.TypeOf(s1Str).String(), s1)

	var s2 Student

	json.Unmarshal(s1Str, &s2)
	fmt.Printf("s2: %v\n", s2)
}
