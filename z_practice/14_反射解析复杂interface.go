package main

import (
	"reflect"
	"encoding/json"
	"fmt"
	"errors"
)

type ShopItemTemplate struct {
	ShopId  int
	ShopTag string
}

type ShopItems []ShopItemTemplate

var shopItemsMap map[string]ShopItems

func main() {
	/*mapInt := make(map[string]int)
	setMapInt(&mapInt)
	mapIntStr, _ := json.Marshal(mapInt)
	fmt.Println("mapIntStr:", string(mapIntStr))
	fmt.Println("=============================================")

	mapIntArr := make(map[string][]int)
	getMapValueType(&mapIntArr)
	mapIntArrStr, _ := json.Marshal(mapIntArr)
	fmt.Println("mapIntArrStr:", string(mapIntArrStr))
	fmt.Println("=============================================")

	mapStruct := make(map[string]ShopItemTemplate)
	getMapStructValueType(&mapStruct)
	mapStructStr, _ := json.Marshal(mapStruct)
	fmt.Println("mapStructStr:", string(mapStructStr))
	fmt.Println("=============================================")*/

	mapStructArr := make(map[string][]ShopItemTemplate)
	getMapStructArrValueType(&mapStructArr)
	mapStructArrStr, _ := json.Marshal(mapStructArr)
	fmt.Println("mapStructArr:", string(mapStructArrStr))
}

func getMapStructArrValueType(o interface{}) error {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Map {
		return errors.New("input should be map, not:" + t.Kind().String())
	}

	type4MapValue := reflect.Indirect(reflect.ValueOf(o)).Type().Elem() //获取map值的类型
	if type4MapValue.Kind() != reflect.Slice {
		return errors.New("input map value should be slice, not:" + type4MapValue.Kind().String())
	}
	fmt.Println("type4MapValue:", type4MapValue)

	type4VauleElem := reflect.Indirect(reflect.New(type4MapValue).Elem()).Type().Elem() //获取map值的元素类型(map值为数组)
	fmt.Println("type4VauleElem:", type4VauleElem)

	ptr := reflect.New(type4VauleElem).Interface()
	//val := reflect.Indirect(reflect.ValueOf(ptr))
	//field, _ := val.Type().FieldByName("ShopId") 	//获取struct结构中字段shop的类型

	s1 := ShopItemTemplate{123, "tag"}
	sStr, _ := json.Marshal(s1)
	json.Unmarshal(sStr, ptr)

	//mapValue := reflect.New(type4MapValue)                                  //声明一个map值的变量
	//mapValue = reflect.Append(mapValue.Elem(), reflect.ValueOf(ptr).Elem()) //将值追加到数组
	//mapValue = reflect.Append(mapValue, reflect.ValueOf(ptr).Elem()) //将值追加到数组
	//reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("aaaa"), mapValue)

	//mapValueEn := reflect.New(type4MapValue) //声明一个map值的变量
	//reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("bbbb"), mapValueEn.Elem())
	//reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("bbbb"), reflect.Append(reflect.ValueOf(o).Elem().MapIndex(reflect.ValueOf("bbbb")), reflect.ValueOf(ptr).Elem()))

	reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("data"), reflect.New(type4MapValue).Elem())
	reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("data_zh_cn"), reflect.New(type4MapValue).Elem())
	reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("data_zh_tw"), reflect.New(type4MapValue).Elem())

	p := reflect.New(type4VauleElem).Interface()
	v := reflect.ValueOf(p)
	vv := reflect.Indirect(v)
	fmt.Println(v.Type().Kind(), vv.Type().Kind())
	//v.Field(0).SetInt(100)
	vv.Field(0).SetInt(10000)
	fmt.Println(v, vv)

	return nil
}

func getMapStructValueType(o interface{}) {
	sl := reflect.Indirect(reflect.ValueOf(o))
	typeOfT := sl.Type().Elem()
	fmt.Println(typeOfT)

	ptr := reflect.New(typeOfT).Interface()
	val := reflect.Indirect(reflect.ValueOf(ptr))
	field, _ := val.Type().FieldByName("ShopId")
	fmt.Println(field.Type.Kind())

	field, _ = val.Type().FieldByName("ShopTag")
	fmt.Println(field.Type.Kind())

	s := ShopItemTemplate{123, "tag"}
	sStr, _ := json.Marshal(s)
	//ptr := reflect.New(typeOfT).Interface()
	json.Unmarshal(sStr, ptr)
	fmt.Println(reflect.ValueOf(ptr).Elem())
	reflect.ValueOf(o).Elem().SetMapIndex(reflect.ValueOf("aaaa"), reflect.ValueOf(ptr).Elem())
}

func getMapValueType(o interface{}) {
	sl := reflect.Indirect(reflect.ValueOf(o))
	typeOfT := sl.Type().Elem()
	fmt.Println(typeOfT)

	val := reflect.New(typeOfT)

	a := reflect.Indirect(val.Elem())
	fmt.Println(a.Type().Elem())

	val = reflect.Append(val.Elem(), reflect.ValueOf(10), reflect.ValueOf(20))
	fmt.Println(val)
	v := reflect.ValueOf(o).Elem()
	v.SetMapIndex(reflect.ValueOf("test"), val)

	//ptr := reflect.New(typeOfT).Interface()
	//val := reflect.Indirect(reflect.ValueOf(ptr))
	//val = reflect.Append(val, reflect.ValueOf(10), reflect.ValueOf(20))
	//fmt.Println(val)

	//v := reflect.ValueOf(o).Elem()
	//fmt.Println(v)
	//v.SetMapIndex(reflect.ValueOf("test"), sl)
}

func setMapInt(o interface{}) {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	v.SetMapIndex(reflect.ValueOf("a"), reflect.ValueOf(1))
	v.SetMapIndex(reflect.ValueOf("b"), reflect.ValueOf(2))
	v.SetMapIndex(reflect.ValueOf("c"), reflect.ValueOf(3))
	for _, idx := range v.MapKeys() {
		fmt.Printf("key: %v, value: %v\n", idx, v.MapIndex(idx))
	}
}
