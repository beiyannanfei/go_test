package main

import (
	"fmt"
)

type littleGirl struct {
	Name string
	Age  int
}

type girl *littleGirl

//todo Go语言中规定，只有类型（Type）和指向他们的指针（*Type）才是可能会出现在接收器声明里的两种接收器，为了避免歧义，明确规定，如果一个类型名本身就是一个指针的话，是不允许出现在接收器中的。
func (this girl) changeName(name string) {	//invalid receiver type girl (girl is a pointer type)
	this.Name = name
}

func main() {
	littleGirl := girl{Name: "Rose", Age: 1}

	girl.changeName("yoyo")
	fmt.Println(littleGirl)
}
