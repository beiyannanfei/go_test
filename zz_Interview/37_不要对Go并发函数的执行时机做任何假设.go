package main

import (
	"fmt"
	"runtime"
	"time"
)

func t1() {
	names := []string{"lily", "yoyo", "cersei", "rose", "annei"}
	for _, name := range names {
		go func() {
			fmt.Println(name)
		}()
	}
	runtime.GOMAXPROCS(1)
	runtime.Gosched()

	/*
	annei
	annei
	annei
	annei
	annei
	*/

	//todo 输出的都是“annei”，而“annei”又是“names”的最后一个元素，那么也就是说程序打印出了最后一个元素的值，而name对于匿名函数来讲又是一个外部的值。因此，我们可以做一个推断：虽然每次循环都启用了一个协程，但是这些协程都是引用了外部的变量，当协程创建完毕，再执行打印动作的时候，name的值已经不知道变为啥了，因为主函数协程也在跑，大家并行，但是在此由于names数组长度太小，当协程创建完毕后，主函数循环早已结束，所以，打印出来的都是遍历的names最后的那一个元素“annei”。
}

func t2() {
	names := []string{"lily", "yoyo", "cersei", "rose", "annei"}
	for _, name := range names {
		go func() {
			fmt.Println(name)
		}()
		time.Sleep(time.Millisecond)
	}
	runtime.GOMAXPROCS(1)
	runtime.Gosched()

	/*
	lily
	yoyo
	cersei
	rose
	annei
	*/

}

func t3() {
	names := []string{"lily", "yoyo", "cersei", "rose", "annei"}
	for _, name := range names {
		temp := name
		go func() {
			fmt.Println(temp)
		}()
	}
	runtime.GOMAXPROCS(1)
	runtime.Gosched()
}

func main() {
	t1()
	fmt.Println("--------------------------------------")
	time.Sleep(time.Millisecond * 100)

	t2()
	fmt.Println("--------------------------------------")
	time.Sleep(time.Millisecond * 100)

	t3()
}

//todo 以上我们得出一个结论，不要对“go函数”的执行时机做任何的假设，除非你确实能做出让这种假设成为绝对事实的保证。