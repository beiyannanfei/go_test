package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	jsonStr := `{"id":1058,"name":"RyuGou"}`
	var jsonData map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &jsonData)

	fmt.Println("type: ", reflect.TypeOf(jsonData["id"])) //type:  float64
	//todo 使用 Golang 解析 JSON  格式数据时，若以 interface{} 接收数据，则会按照下列规则进行解析：
	//bool, for JSON booleans
	//float64, for JSON numbers
	//string, for JSON strings
	//[]interface{}, for JSON arrays
	//map[string]interface{}, for JSON objects
	//nil for JSON null
}
