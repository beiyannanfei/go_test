package main

import (
	"time"
	"fmt"
)

func main() {
	t := time.Now().Format(time.RFC3339)
	fmt.Println(t)
	t = time.Now().Format("20060102150405")
	fmt.Println(t)
}
