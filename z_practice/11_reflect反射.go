package main

import (
	"reflect"
	"fmt"
)

type X int
type Y int

func typeCheck() {
	var a X = 100
	var b X = 200
	var c Y = 300
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	tc := reflect.TypeOf(c)
	//判断真实类型是否相等
	fmt.Printf("ta == tb -> %v, ta == tc -> %v\n", ta == tb, ta == tc) //ta == tb -> true, ta == tc -> false
	//判断底层类型是否相等
	fmt.Printf("ta.Kind == tc.Kind -> %v\n", ta.Kind() == tb.Kind()) //ta.Kind == tc.Kind -> true
}

func buildCompositeType() {
	arrType := reflect.ArrayOf(10, reflect.TypeOf(byte(0)))
	mapType := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(0))
	fmt.Printf("buildCompositeType arrType: %v, mapType: %v\n", arrType, mapType) //buildCompositeType arrType: [10]uint8, mapType: map[string]int
}

func diffPtr() {
	x := 100
	tx := reflect.TypeOf(x)
	tp := reflect.TypeOf(&x)
	fmt.Printf("tx: %v, tp: %v, tx==tp -> %v\n", tx, tp, tx == tp)                         //tx: int, tp: *int, tx==tp -> false
	fmt.Printf("tx.Kind: %v, tp.Kind: %v\n", tx.Kind(), tp.Kind())                         //tx.Kind: int, tp.Kind: ptr
	fmt.Printf("tx: %v, tp.Elem: %v, tx==tp.Elem -> %v\n", tx, tp.Elem(), tx == tp.Elem()) //tx: int, tp.Elem: int, tx==tp.Elem -> true
}

type user struct {
	name string
	age  int
}

type manager struct {
	user
	title string
}

func reflectStruct() { //对含有匿名字段的结构体进行反射遍历
	var m manager
	tp := reflect.TypeOf(&m)
	fmt.Printf("reflectStruct t: %v, t.Kind: %v, t.Elem: %v\n", tp, tp.Kind(), tp.Elem()) //reflectStruct t: *main.manager, t.Kind: ptr, t.Elem: main.manager
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		fmt.Printf("Name: %v, Type: %v, Offset: %v\n", f.Name, f.Type, f.Offset)
		if f.Anonymous { //输出匿名字段结构
			for x := 0; x < f.Type.NumField(); x++ {
				af := f.Type.Field(x)
				fmt.Printf("    Name: %4v, Type: %6v, Offset: %v\n", af.Name, af.Type, af.Offset)
			}
		}
	}
	/*
	Name: user, Type: main.user, Offset: 0
    	Name: name, Type: string, Offset: 0
    	Name:  age, Type:    int, Offset: 16
	Name: title, Type: string, Offset: 24
	*/
}

type A struct {
	a int
}
type B struct {
	A
	b int
}

func (a A) Av() {
	fmt.Println("***** Av")
}
func (a *A) Ap() {
	fmt.Println("***** Ap")
}
func (b B) Bv() {
	fmt.Println("***** Bv")
}

func (b *B) Bp() {
	fmt.Println("***** Bp")
}

func funcType() {
	var b B
	t := reflect.TypeOf(&b)
	ts := []reflect.Type{t, t.Elem()}
	for index, tValue := range ts {
		fmt.Printf("index: %v, tValue: %v :\n", index, tValue)
		for i := 0; i < tValue.NumMethod(); i++ {
			fmt.Println("	", tValue.Method(i))
		}
	}
	/*
	index: 0, tValue: *main.B :
	 	{Ap  func(*main.B) <func(*main.B) Value> 0}
	 	{Av  func(*main.B) <func(*main.B) Value> 1}
	 	{Bp  func(*main.B) <func(*main.B) Value> 2}
	 	{Bv  func(*main.B) <func(*main.B) Value> 3}
	index: 1, tValue: main.B :
		{Av  func(main.B) <func(main.B) Value> 0}
	 	{Bv  func(main.B) <func(main.B) Value> 1}
	*/
}

func main() {
	fmt.Println("---------- 在面对类型时，需要区分 Type 和 Kind。前者表示真实类型（静态类型），后者表示其基础结构（底层类型）类别 -- 基类型 ----------")

	var a X = 10
	t := reflect.TypeOf(a)
	fmt.Println("t: ", t)           //t: main.X		Type 表示真实类型（静态类型）
	fmt.Println("kind: ", t.Kind()) //kind:  int		Kind 表示其基础结构（底层类型）类别 -- 基类型
	fmt.Println("name: ", t.Name()) //name:  X

	fmt.Println("----------------------- 在类型判断上，须选择正确的方式 -----------------------")
	typeCheck()

	fmt.Println("----------------------- 除通过实际对象获取类型外，也可直接构造一些基础复合类型 -----------------------")
	buildCompositeType()

	fmt.Println("----------------- 传入对象 应区分 基类型 和 指针类型，因为它们并不属于同一类型 -----------------")
	diffPtr()

	fmt.Println("----------------- 方法 Elem() 返回 指针、数组、切片、字典（值）或 通道的 基类型 -----------------")
	fmt.Printf("map[string]int{}基类型: %v\n", reflect.TypeOf(map[string]int{}).Elem()) //map[string]int{}基类型: int
	fmt.Printf("[]int32{}基类型: %v\n", reflect.TypeOf([]int32{}).Elem())               //[]int32{}基类型: int32

	fmt.Println("----------------- 只有在获取 结构体指针 的 基类型 后，才能遍历它的字段 -----------------")
	reflectStruct() //对含有匿名字段的结构体进行反射遍历

	fmt.Println("----------------- 对于匿名字段，可用多级索引（按照定义顺序）直接访问 -----------------")
	var m = manager{user{"bynf", 29}, "engineer"}
	t = reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	n, find := t.FieldByName("name")                                                         //按名称查找
	fmt.Printf("n: %v, find: %v\n", n, find)                                                 //n: {name main string  0 [0 0] false}, find: true
	fmt.Printf("n.Name: %v, n.Type: %v, value: %v\n", n.Name, n.Type, v.FieldByName("name")) //n.Name: name, n.Type: string, value: bynf

	ageT := t.FieldByIndex([]int{0, 1}) // 按多级索引查找
	ageV := v.FieldByIndex([]int{0, 1})
	fmt.Printf("ageT: %v, ageV: %v\n", ageT, ageV)                                            //ageT: {age main int  16 [1] false}, ageV: 29
	fmt.Printf("ageT.Name: %v, ageT.Type: %v, value: %v\n", ageT.Name, ageT.Type, ageV.Int()) //ageT.Name: age, ageT.Type: int, value: 29

	fmt.Println("----------------- 输出方法集时，一样区分 基类型 和 指针类型 -----------------")
	funcType()
}
