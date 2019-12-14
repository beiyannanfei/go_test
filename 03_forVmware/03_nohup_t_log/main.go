package main

import (
	"fmt"
	"os"
	"time"
)

func showlog() {
	fmt.Printf("[%s]================info\n", time.Now().Format("2006-01-02 15:04:05"))
}

func showerr() {
	fmt.Fprintf(os.Stderr, "[%s]====================err\n", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	for {
		showlog()
		showerr()
		time.Sleep(time.Second)
	}
}
