package main

import (
	"path"
	"fmt"
	"strings"
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
	pathStr := "/10113 "
	fmt.Println(path.Base(pathStr))
}

func main() {
	isAbs()
	split()
	filename := "templates/achievement.xlsx"
	fmt.Println(path.Base(filename))
	fmt.Println(path.Ext(filename))
	fmt.Println(strings.TrimRight(path.Base(filename), path.Ext(filename)))
}
