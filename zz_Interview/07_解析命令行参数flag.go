package main

import (
	"flag"
	"fmt"
)

//https://studygolang.com/articles/17368

var name string
var addr string

func init() {
	flag.StringVar(&name, "name", "nobody", "请设置名字")
	flag.StringVar(&addr, "addr", "default", "请设置地址")
}

func main() {
	flag.Parse()
	fmt.Println("Hello,", name)
	fmt.Println("addr,", addr)
}

//go run 07_解析命令行参数.go -name a -addr b
