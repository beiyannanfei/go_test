package main

import (
	"time"
	"math/rand"
	"fmt"
)

//https://blog.csdn.net/skh2015java/article/details/60330975

//上面的例子都使用一个信道作为返回值，可以把信道的数据合并到一个信道的。 不过这样的话，我们需要按顺序输出我们的返回值（先进先出）。
//如下，我们假设要计算很复杂的一个运算 100-x , 分为三路计算， 最后统一在一个信道中取出结果:

func do_stuff(x int) int { //一个比较耗时的事情，比如计算
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond) //延时,模拟计算
	return 100 - x
}

func branch(x int) chan int { //每个分支开出一个goroutine做计算并把计算结果流入各自信道
	ch := make(chan int)

	go func() {
		ch <- do_stuff(x)
	}()

	return ch
}

func fanIn(chs ... chan int) chan int {
	ch := make(chan int)

	for _, c := range chs {
		go func(c chan int) { //复合 注意此处明确传值
			ch <- <-c
		}(c)
	}

	return ch
}

func main() {
	results := fanIn(branch(1), branch(2), branch(3))

	for i := 0; i < 3; i++ {
		v, ok := <-results
		fmt.Println(v, ok)
	}
}
