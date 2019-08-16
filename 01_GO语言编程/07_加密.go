package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func main() {
	testString := "hello world"
	Md5Inst := md5.New()
	Md5Inst.Write([]byte(testString))
	Resutl := Md5Inst.Sum([]byte(""))
	fmt.Printf("%x\n\n", Resutl)

	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(testString))
	Resutl = Sha1Inst.Sum([]byte(""))
	fmt.Printf("%x\n\n", Resutl)
}
