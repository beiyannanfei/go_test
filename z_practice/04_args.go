package main

import (
	"os"
	"fmt"
)

//命令行参数
// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 04_args.go aaa bbb ccc
func main() {
	argList := os.Args
	fmt.Printf("argList: %#v\n", argList)		//argList: []string{"/var/folders/q1/d39w_8sx5t30j86vzn2s3xcc0000gn/T/go-build903566758/b001/exe/04_args", "aaa", "bbb", "ccc"}

	for i, arg := range argList {
		fmt.Printf("args[%v] = %s\n", i, arg)
	}
}
