package main

import (
	"fmt"
	"crypto/md5"
)
// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 01_md5.go
func main() {
	var base = fmt.Sprintf("%d-%s-%d", 12, "abcd", 34)
	fmt.Printf("base: %v\n, md5Value: %x\n", base, md5.Sum([]byte(base)))
}
