package main

import (
	"time"
	"fmt"
)

func main() {
	go func() {
		for {
			f5()
			now := time.Now()
			next := now.Add(time.Second * 5)
			t := time.NewTimer(next.Sub(now))
			o := <-t.C
			fmt.Printf("======== %v\n", o)
		}
	}()

	time.Sleep(100 * time.Second)

	/*now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05")
	fmt.Println("nowStr: ", nowStr)

	next := now.Add(time.Hour * 24)
	nextStr := next.Format("2006-01-02 15:04:05")
	fmt.Println("nextStr: ", nextStr)

	fmt.Println(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), next.Location())*/
}

func f5() {
	fmt.Printf("[%v] 每5s输出一次.\n", time.Now().Format("2006-01-02 15:04:05"))
}
