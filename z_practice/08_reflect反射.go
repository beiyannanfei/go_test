package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Id   int
	Name string
	Age  int
}

func (u User) SayHello() {
	fmt.Println("Hello", u.Name)
}

func (u User) HiUser(name string) {
	fmt.Printf("Hi %v, User Name: %v\n", name, u.Name)
}

func Info(o interface{}) {
	t := reflect.TypeOf(o)                 //反射使用 TypeOf 和 ValueOf 函数从接口中获取目标对象信息
	fmt.Println("param o Type:", t.Name()) // param o Type: User

	v := reflect.ValueOf(o)
	fmt.Printf("%#v\n", v) // main.User{Id:123, Name:"bynf", Age:29}
	fmt.Println(v.FieldByName("Name"))	//bynf

	//通过索引来取得它的所有字段，这里通过t.NumField来获取它多拥有的字段数量，同时来决定循环的次数
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)             //通过这个i作为它的索引，从0开始来取得它的字段
		value := v.Field(i).Interface() //通过interface方法来取出这个字段所对应的值
		fmt.Printf("field: (%v)%v, value: %v\n", field.Type, field.Name, value)
	}

	//通过t.NumMethod来获取它拥有的方法的数量，来决定循环的次数
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fmt.Printf("method type: %v, name: %v\n", method.Type, method.Name) //method type: func(main.User), name: SayHello
	}

	kind := t.Kind() //通过kind方法判断传入的类型是否是我们需要反射的类型
	fmt.Printf("param kind: %v, is reflect.Struct: %v\n", kind, kind == reflect.Struct)
}

type Manager struct {
	User //反射会将匿名字段作为一个独立字段来处理
	Title string
}

func Set(o interface{}) {
	v := reflect.ValueOf(o)
	fmt.Printf("Set v: %#v\n", v) //Set v: &main.User{Id:789, Name:"Tom", Age:25}
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		fmt.Println("v Can not set")
		return
	}

	fmt.Printf("Set v.Elem: %#v\n", v.Elem()) //Set v.Elem: main.User{Id:789, Name:"Tom", Age:25}
	f := v.Elem().FieldByName("Name")
	fmt.Printf("f: %#v\n", f) //f: "Tom"
	if !f.IsValid() {
		fmt.Println("修改失败")
	}

	if f.Kind() == reflect.String {
		f.SetString("beiyannanfei")
	}
}

func main() {
	fmt.Println("----------------------------对某一个struct进行反射的基本操作---------------------------")

	u := User{123, "bynf", 29}
	u.SayHello()
	Info(u)

	fmt.Println("----------------------------反射 匿名或嵌入字段---------------------------")

	m := Manager{User: User{1112233, "Jack", 12}, Title: "AABBCC"}
	t := reflect.TypeOf(m)
	//m Name is: Manager, Type is: struct
	fmt.Printf("m Name is: %v, Type is: %v\n", t.Name(), t.Kind())
	// #号会将reflect的struct的详情页打印出来，可以看出来这是一个匿名字段
	//m field(0): reflect.StructField{Name:"User", PkgPath:"", Type:(*reflect.rtype)(0x10d4660), Tag:"", Offset:0x0, Index:[]int{0}, Anonymous:true}
	fmt.Printf("m field(0): %#v\n", t.Field(0))
	v := reflect.ValueOf(m)
	//m reflect.ValueOf: main.Manager{User:main.User{Id:1112233, Name:"Jack", Age:12}, Title:"AABBCC"}
	fmt.Printf("m reflect.ValueOf: %#v\n", v)
	//m field(0) Name: User, Type: main.User, Value: {1112233 Jack 12}
	fmt.Printf("m field(0) Name: %v, Type: %v, Value: %v\n", t.Field(0).Name, t.Field(0).Type, v.Field(0).Interface())

	//我们就可以将User当中的ID取出来,这里面需要传进方法中的是一个int类型的slice，User相对于manager索引是0，id相对于User索引也是0
	mUserIdKey := t.FieldByIndex([]int{0, 0})
	//mUserIdKey: reflect.StructField{Name:"Id", PkgPath:"", Type:(*reflect.rtype)(0x10bf360), Tag:"", Offset:0x0, Index:[]int{0}, Anonymous:false}, Name: Id, Type: int
	fmt.Printf("mUserIdKey: %#v, Name: %v, Type: %v\n", mUserIdKey, mUserIdKey.Name, mUserIdKey.Type)

	//同理，可以将User中Id的值取出来
	mUserIdValue := v.FieldByIndex([]int{0, 0})
	//mUserIdValue: 1112233
	fmt.Printf("mUserIdValue: %#v\n", mUserIdValue)

	fmt.Println("----------------------------通过反射修改struct中的内容---------------------------")

	x := 123
	v = reflect.ValueOf(&x) //传递指针才能修改
	v.Elem().SetInt(456)
	fmt.Printf("modify x: %v, Type: %v, Value: %v\n\n", x, v.Type(), v.Elem()) //modify x: 456, Type: *int, Value: 456

	u = User{789, "Tom", 25}
	Set(&u)
	fmt.Printf("modify u: %#v\n", u) //modify u: main.User{Id:789, Name:"beiyannanfei", Age:25}

	fmt.Println("----------------------------通过反射进行方法的调用 动态调用方法---------------------------")

	u = User{258, "BYNF", 12}
	v = reflect.ValueOf(u)
	mv := v.MethodByName("HiUser")
	args := []reflect.Value{reflect.ValueOf("JOE")}
	mv.Call(args) // Hi JOE, User Name: BYNF
}
