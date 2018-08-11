package main

import (
	"reflect"
	"fmt"
	"strings"
)

//反射库提供了内置函数make和new的对应操作，如reflect.MakeFunc()方法，通过它可以实现类似“泛型”的功能(https://www.cnblogs.com/susufufu/archive/2017/10/12/7653579.html)

func add(args []reflect.Value) (results []reflect.Value) {
	if len(args) <= 0 {
		return nil
	}

	var r reflect.Value
	switch args[0].Kind() {
	case reflect.Int:
		sum := 0
		for _, v := range args {
			fmt.Printf("reflect.Int For v(%v): %v\n", v.Kind(), v) //reflect.Int For v(int): 12		reflect.Int For v(int): 23
			sum += int(v.Int())
		}
		r = reflect.ValueOf(sum)
	case reflect.String:
		sumStr := make([]string, 0, len(args))
		for _, v := range args {
			fmt.Printf("reflect.String For v(%v): %#v\n", v.Kind(), v) //reflect.String For v(string): "Hello"			reflect.String For v(string): "world"
			sumStr = append(sumStr, v.String())
		}
		r = reflect.ValueOf(strings.Join(sumStr, " "))
	}

	results = append(results, r)
	return
}

func makeAdd(T interface{}) {
	value := reflect.ValueOf(T)
	fmt.Printf("value: %#v\n", value) //value: (*func(int, int) int)(0xc42000c028)
	fn := value.Elem()
	fmt.Printf("fn: %#v\n", fn) //fn: (func(int, int) int)(nil)

	v := reflect.MakeFunc(fn.Type(), add)    //把原始函数变量的类型和通用算法函数存到同一个Value中
	fmt.Printf("v: %v\n", reflect.TypeOf(v)) //v: reflect.Value
	fn.Set(v)                                //把原始函数指针变量指向v，这样它就获得了函数体
}

func main() {
	var intAdd func(x int, y int) int
	makeAdd(&intAdd)
	sumInt := intAdd(12, 23)
	fmt.Printf("sumInt: %#v\n", sumInt) //sumInt: 35
	fmt.Println("----------------------------------------------------")

	var strAdd func(x string, y string) string
	makeAdd(&strAdd)
	sumStr := strAdd("Hello", "world")
	fmt.Printf("sumStr: %#v\n", sumStr) //sumStr: "Hello world"
}
