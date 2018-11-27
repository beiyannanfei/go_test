package main

import (
	"path/filepath"
	"os"
	"fmt"
	"path"
	"strings"
)

func main() {
	dir := "./"
	err := filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if path.Ext(filename) != ".go" {
			return nil
		}

		if strings.HasPrefix(filename, "vendor") {
			return nil
		}

		fmt.Println(filename)
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
