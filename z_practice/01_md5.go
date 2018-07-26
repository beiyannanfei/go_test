package main

import (
	"fmt"
	"crypto/md5"
)

func main() {
	var base = fmt.Sprintf("%d-%s-%d", 12, "abcd", 34)
	fmt.Printf("base: %v\n, md5Value: %x\n", base, md5.Sum([]byte(base)))
}
