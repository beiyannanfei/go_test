package main

import (
	"fmt"
	"github.com/beiyannanfei/go_test/qsc_test/002_读写锁/lock"
	"time"
)

func main() {
	go lock.BuildTree()

	time.Sleep(2 * time.Second)
	s := lock.GetSubTreeByDepartmentId(10)
	fmt.Println(s)
	select {}
}
