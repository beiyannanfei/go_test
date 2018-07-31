package main

import (
	"io/ioutil"
	"fmt"
)
// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 05_readFile.go
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
