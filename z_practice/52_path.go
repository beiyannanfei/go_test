package main

import (
	"path"
	"fmt"
)

func isAbs() {
	//IsAbs返回路径是否是一个绝对路径。
	dirPath := "/dev/null"
	r := path.IsAbs(dirPath)
	fmt.Printf("%v IsAbs: %v\n", dirPath, r)

	dirPath = "../a.txt"
	r = path.IsAbs(dirPath)
	fmt.Printf("%v IsAbs: %v\n", dirPath, r)
}

func split() {
	pathStr := "a/b/c/d.txt"

}

func main() {
	isAbs()
}
