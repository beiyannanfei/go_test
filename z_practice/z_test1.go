package main

import (
	"strings"
	"fmt"
)

func main() {
	fmt.Println(strings.Replace(strings.TrimSpace(""), "http", "https", 1))
}
