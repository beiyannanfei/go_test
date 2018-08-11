package main

import (
	"fmt"
	"reflect"
)

func main() {
	fmt.Println("-------- 传入指针，一样需要通过 Elem() 获取目标对象。因为被接口存储的指针本身是不能寻址和进行设置操作的 --------")
	a := 100
	va := reflect.ValueOf(a)
	vp := reflect.ValueOf(&a).Elem()
	fmt.Printf("va.CanAddr: %v, va.CanSet: %v\n", va.CanAddr(), va.CanSet()) //va.CanAddr: false, va.CanSet: false
	fmt.Printf("vp.CanAddr: %v, vp.CanSet: %v\n", vp.CanAddr(), vp.CanSet()) //vp.CanAddr: true, vp.CanSet: true
	vp.Set(reflect.ValueOf(123))
	fmt.Printf("final a: %v\n", a) //final a: 123

	fmt.Println("-------- 复合类型对象设置 --------")
	c := make(chan int, 4)
	vc := reflect.ValueOf(c)
	if vc.TrySend(reflect.ValueOf(100)) {
		vcRecv, vcBool := vc.TryRecv()
		fmt.Printf("vcRecv: %v, vcBool: %v\n", vcRecv, vcBool) //vcRecv: 100, vcBool: true
	}

	fmt.Println("-------- 接口有两种 nil 状态，这一直是个潜在麻烦。解决方法是用 IsNil() 判断值是否为 nil --------")
	var ia interface{} = nil
	var ib interface{} = (*int)(nil)
	fmt.Printf("ia is null: %v\n", ia == nil)                 //ia is null: true
	fmt.Printf("ib is null: %v\n", ib == nil)                 //ib is null: false
	fmt.Printf("ib isNil: %v\n", reflect.ValueOf(ib).IsNil()) //ib isNil: true

	fmt.Println("-------- Value 里的某些方法并未实现 ok-idom 或返回 error，所以得自行判断返回的是否为 Zero Value --------")
	vs := reflect.ValueOf(
		struct {
			name string
			age  int
		}{},
	)
	fmt.Printf("field 'name' isValid: %v\n", vs.FieldByName("name").IsValid())       //field 'name' isValid: true
	fmt.Printf("field 'age' isValid: %v\n", vs.FieldByName("age").IsValid())         //field 'age' isValid: true
	fmt.Printf("field 'address' isValid: %v\n", vs.FieldByName("address").IsValid()) //field 'address' isValid: false
}
