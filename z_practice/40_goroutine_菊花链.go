package main

import (
	"fmt"
	"time"
)

func xrange_40() chan int { //从2开始自增的整数生成器
	var ch = make(chan int)

	go func() {
		for i := 2; ; i++ {
			ch <- i //直到信道索要数据，才把i添加进信道
		}
	}()

	return ch
}

func filter(in chan int, number int) chan int {
	// 输入一个整数队列，筛出是number倍数的, 不是number的倍数的放入输出队列
	// in:  输入队列
	out := make(chan int)

	go func() {
		for {
			i := <-in //从输入中区一个
			fmt.Printf("filter i: %v, number: %v\n", i, number)
			if i%number != 0 {
				out <- i
			}
		}
	}()

	return out
}

func main() {
	const max = 100     //找出100以内的所有素数
	nums := xrange_40() //初始化一个整数生成器
	number := <-nums    //从生成器中抓一个整数(2), 作为初始化整数

	for number <= max { //number作为筛子，当筛子超过max的时候结束筛选
		fmt.Println("number:", number) //打印素数, 筛子即一个素数
		nums = filter(nums, number)    //筛掉number的倍数
		number = <-nums                //更新筛子
		time.Sleep(time.Second)
	}
}
