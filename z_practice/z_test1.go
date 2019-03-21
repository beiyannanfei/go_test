package main

import (
	"fmt"
	"time"
)

func main() {
	str := "asdf"
	fmt.Println(len(str))

	fmt.Println(time.Now().Unix() - 30*24*3600)
}
