package main

import (
	"fmt"
	"reflect"
)

//https://blog.csdn.net/liumiaocn/article/details/55253632
func main() {
	fmt.Println(`---------------------- “接口类型变量”=>“反射类型对象” ----------------------`)

	var circle1 = 6.28
	var icir interface{}

	icir = circle1
	fmt.Printf("Reflect: icir.Value: %v\n", reflect.ValueOf(icir)) //Reflect: icir.Value: 6.28
	fmt.Printf("Reflect: icir.Type: %v\n", reflect.TypeOf(icir))   //Reflect: icir.Type: float64

	//可以看到ValueOf和TypeOf的参数都是空接口，因此，这说明可以直接使用变量传进去
	fmt.Printf("Reflect: circle1.Value: %v\n", reflect.ValueOf(circle1)) //Reflect: circle1.Value: 6.28
	fmt.Printf("Reflect: circle1.Type: %v\n", reflect.TypeOf(circle1))   //Reflect: circle1.Type: float64

	fmt.Println(`---------------------- “反射类型对象”=>“接口类型变量” ----------------------`)

	valueRef := reflect.ValueOf(icir)
	fmt.Printf("valueRef: %v\n", valueRef)                       //valueRef: 6.28
	fmt.Printf("valueRef Interface: %v\n", valueRef.Interface()) //valueRef Interface: 6.28
	y := valueRef.Interface().(float64)
	fmt.Printf("y: %v, type: %v\n", y, reflect.TypeOf(y)) //y: 6.28, type: float64

	fmt.Println(`---------------------- 修改“反射类型对象” ----------------------`)

	value := reflect.ValueOf(circle1)
	fmt.Printf("reflect value: %v\n", value) //reflect value: 6.28
	//value.SetFloat(3.14) error value变量之所以是不可写的，因为其所指向的是一个副本，因此不具有可写性
	//go中和特意提供了一个CanSet函数可以进行确认是否是settable的
	fmt.Printf("reflect value can set: %v\n", value.CanSet()) //reflect value can set: false

	value2 := reflect.ValueOf(&circle1)
	fmt.Printf("reflect value2 can set: %v\n", value2.CanSet()) //reflect value2 can set: false

	value3 := value2.Elem()
	fmt.Printf("value3: %v, reflect value3 can set: %v\n", value3, value3.CanSet()) //value3: 6.28, reflect value3 can set: true
	value3.SetFloat(3.14)
	fmt.Printf("final value3: %v, circle1: %v\n", value3, circle1) //final value3: 3.14, circle1: 3.14
	value3.Set(reflect.ValueOf(1.234))
	fmt.Printf("final value3: %v, circle1: %v\n", value3, circle1) //final value3: 1.234, circle1: 1.234

	fmt.Println(`---------------------- 修改“反射类型对象”-结构体 ----------------------`)
	type T struct {
		A int
		B string
	}

	t := T{369, "AABBCCDD"}
	v := reflect.ValueOf(&t)
	fmt.Printf("v: %#v\n", v) //v: &main.T{A:369, B:"AABBCCDD"}
	s := v.Elem()
	fmt.Printf("s: %#v\n", s) //s: main.T{A:369, B:"AABBCCDD"}
	typeOfT := s.Type()
	fmt.Printf("typeOfT: %v\n", typeOfT) //typeOfT: main.T
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i) //迭代s的各个域，注意每个域仍然是反射。
		fmt.Printf("i: %v, Name: %v, Type: %v, value: %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		//i: 0, Name: A, Type: int, value: 369
		//i: 1, Name: B, Type: string, value: AABBCCDD
	}

	//两种方式修改结构体内容
	s.FieldByName("A").Set(reflect.ValueOf(123))	//FieldByName通过字段名获取变量，Set方法设置变量，通过reflect.ValueOf反射获取变量类型
	s.Field(1).SetString("asdf")	//Field通过索引获取变量
	fmt.Printf("t now is: %#v\n", t) //t now is: main.T{A:123, B:"asdf"}
}
