package main

import "fmt"

func main() {
	defer_call()
	fmt.Println("---------------------------------------")
	defer_call1()
	fmt.Println("---------------------------------------")
	defer_call2()
}

func defer_call2() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("3", err)
		}

		fmt.Println("打印前")
	}()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("2", err)
		}

		fmt.Println("打印中")
	}()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("1", err)
		}

		fmt.Println("打印后")
	}()

	panic("defer_call2 exception")
}

func defer_call1() {
	defer func() {
		fmt.Println("打印前")
	}()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}

		fmt.Println("打印中")
	}()

	defer func() {
		fmt.Println("打印后")
	}()

	panic("defer_call1 exception")
}

func defer_call() {
	defer func() {
		fmt.Println("before log")
	}()

	defer func() {
		fmt.Println("logging...")
	}()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}

		fmt.Println("after log")
	}()

	panic("defer_call exception")
}

/*
defer_call exception
after log
logging...
before log
---------------------------------------
打印后
defer_call1 exception
打印中
打印前
---------------------------------------
1 defer_call2 exception
打印后
打印中
打印前
*/
