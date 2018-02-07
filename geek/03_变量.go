package main

import (
	"os"
	"fmt"
	"runtime"
)

func main() {
	pkgOs()
	pkgRuntime()
}

func pkgOs() { //os包
	var (
		HOME   = os.Getenv("HOME")
		USER   = os.Getenv("USER")
		GOROOT = os.Getenv("GOROOT")
		path   = os.Getenv("PATH")
	)
	fmt.Printf("HOME: %v, \nUSER: %v, \nGOROOT: %v, \npath: %v\n", HOME, USER, GOROOT, path)
}

func pkgRuntime() {
	goos := runtime.GOOS
	fmt.Printf("The operating system is: %s\n", goos)
}
