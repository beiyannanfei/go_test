package main

import (
	"fmt"
	"crypto/md5"
	"time"
	"strconv"
)

// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 01_md5.go
func main() {
	var base = fmt.Sprintf("%d-%s-%d", 12, "abcd", 34)
	fmt.Printf("base: %v\n, md5Value: %x\n", base, md5.Sum([]byte(base)))

	var cruTime = time.Now().Unix()
	fmt.Printf("cruTime: %v, md5Vaule: %x\n", cruTime, md5.Sum([]byte(strconv.FormatInt(cruTime, 10))))
}
