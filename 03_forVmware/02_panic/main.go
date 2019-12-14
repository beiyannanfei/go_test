package main

import (
	"github.com/beiyannanfei/go_test/03_forVmware/02_panic/t"
	"time"
)

func main() {
	for {
		t.Show()
		t.Add()
		time.Sleep(time.Second)
	}
}
