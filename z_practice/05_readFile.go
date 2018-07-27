package main

import (
	"io/ioutil"
	"fmt"
)

func main() {
	buff, err := ioutil.ReadFile("./05_readFile.go")
	if err != nil {
		fmt.Printf("ioutil.ReadFile err: %v\n", err.Error())
		return
	}

	fmt.Printf("buff: %#v\n", buff)

	fileStr := string(buff)
	fmt.Printf("fileStr: %v\n", fileStr)
}
