package main

import (
	"sync"
	"fmt"
	"time"
)

var a05 string
var once sync.Once

func setup() {
	a05 = "hello world"
	fmt.Println("run setup.")
}

func doProint() {
	once.Do(setup)
	fmt.Println(a05)
}

func main() {
	go doProint()
	go doProint()
	time.Sleep(time.Second)
}
//对于从全局的角度只需要运行一次的代码，比如全局初始化操作，Go语言提供了一个Once类型来保证全局的唯一性操作，具体代码如下:
//如果这段代码没有引入Once，setup()将会被每一个goroutine先调用一次，这至少对于这个例子是多余的。
// 在现实中，我们也经常会遇到这样的情况。Go语言标准库为我们引入了Once类 型以解决这个问题。
// once的Do()方法可以保证在全局范围内只调用指定的函数一次(这里指 setup()函数)，
// 而且所有其他goroutine在调用到此语句时，将会先被阻塞，直至全局唯一的 once.Do()调用结束后才继续。