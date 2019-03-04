package main

import (
	"fmt"
	"time"
)

func sum(s []int, c chan int, sleep int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	time.Sleep(time.Second * time.Duration(sleep))
	c <- sum
}

func main() {
	fmt.Printf("[%v] begin\n", time.Now().Format("2006-01-02 15:04:05"))
	s := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(s[:len(s)/2], c, 1)
	go sum(s[len(s)/2:], c, 2)
	x := <-c
	fmt.Println(x)
	y := <-c
	fmt.Println(y)
	fmt.Printf("[%v] end\n", time.Now().Format("2006-01-02 15:04:05"))
}
