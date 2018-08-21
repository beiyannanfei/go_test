package main

import (
	"sync"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS( /*runtime.NumCPU()*/ 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			fmt.Println("A:", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			fmt.Println("BBBBBBBBBB:", i)
		}
	}()

	wg.Wait()
}
