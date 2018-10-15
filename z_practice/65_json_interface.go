package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		fmt.Println("json Unmarshal err:", err)
		return
	}

	fmt.Println("f:", f)

	m, ok := f.(map[string]interface{})
	if !ok {
		fmt.Println("assert err:", err)
		return
	}

	fmt.Println("m:", m)

	for key, value := range m {
		fmt.Printf("key: %v, value: %v\n", key, value)
		switch vv := value.(type) {
		case string:
			fmt.Printf("key: %v, is a string, value: %v\n", key, value)
		case int:
			fmt.Printf("key: %v, is a int, value: %v\n", key, value)
		case []interface{}:
			fmt.Printf("key: %v, is an array\n", key)
			for i, u := range vv {
				fmt.Printf("%v: %v\n", i, u)
			}
		default:
			fmt.Printf("%v is of a type I don't know how to handler, reflect.TypeOf: %v\n", key, reflect.TypeOf(value))
		}
	}
}
