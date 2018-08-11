package main

import (
	"fmt"
	"reflect"
)

type User12 struct {
	Name string `json:"name" default:"bynf"`
	Age  int    `json:"age" default:"29"`
}

func main() {
	fmt.Println("---------- 反射提取 struct tag，还能自动分解。其常用于 ORM 映射，或数据格式验证 ----------")
	var u User12
	t := reflect.TypeOf(u)
	fmt.Println("t: ", t)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("f.Name: %4v, f.Type: %6v, tag.json: %4v, tag.default: %4v\n", f.Name, f.Type, f.Tag.Get("json"), f.Tag.Get("default"))
	}
	/*
	f.Name: Name, f.Type: string, tag.json: name, tag.default: bynf
	f.Name:  Age, f.Type:    int, tag.json:  age, tag.default:   29
	*/

}
