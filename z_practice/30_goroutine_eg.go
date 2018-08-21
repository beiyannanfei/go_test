package main

import (
	"fmt"
	"reflect"
)

//interface类型chan，取出后转化为对应类型

type user_30 struct {
	Name string
}

func main() {
	userChan := make(chan interface{}, 1)

	u := user_30{Name: "bynf"}
	userChan <- &u
	close(userChan)

	var u1 interface{}
	u1 = <-userChan

	var u2 *user_30
	u2, ok := u1.(*user_30)
	if !ok {
		fmt.Println("can not convert.")
		return
	}
	fmt.Println(u2.Name)
	fmt.Println("------------------------------------")

	uChan := make(chan user_30, 1)
	uu := user_30{Name: "test"}
	uChan <- uu

	uu1, ok := <-uChan
	fmt.Println(uu1, ok, reflect.TypeOf(uu1), uu1.Name) //{test} true main.user_30 test
}
