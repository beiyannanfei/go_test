package main

import (
	"runtime"
	"sync"
	"fmt"
)

func main() {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 10; i++ { //todo i是外部for的一个变量，地址不变化。遍历完成后，最终i=10。 故go func执行时，i的值始终是10。
		go func() {
			fmt.Println("i1: ", i)
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ { //todo go func中i是函数参数，与外部for中的i完全是两个变量。 尾部(i)将发生值拷贝，go func内部指向值拷贝地址
		go func(i int) {
			fmt.Println("i2: ", i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
