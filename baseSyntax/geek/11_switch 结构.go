package main

import "fmt"

func main() {
	fallthrough_test()
	sw1()
	sw2()
	sw3()
}

//一旦成功地匹配到某个分支，在执行完相应代码后就会退出整个 switch 代码块，也就是说您不需要特别使用 break 语句来表示结束。
//因此，程序也不会自动地去执行下一个分支的代码。如果在执行完每个分支的代码后，还希望继续执行后续分支的代码，可以使用 fallthrough 关键字来达到目的。
func fallthrough_test() {
	i := 0
	switch i {
	case 0:
		fallthrough //会继续往下执行
	case 1:
		fmt.Println("=== case 1 ===")
		fallthrough
	case 2:
		fmt.Println("=== case 2 ===") //到此为止
	case 3:
		fmt.Println("=== case 3 ===") //因为上一句代码没有fallthrough，所以不会执行
	default:
		fmt.Println("=== default ===")
	}
}

func sw1() {
	var num1 int = 100

	switch num1 {
	case 98, 99:
		fmt.Println("It's equal to 98")
	case 100:
		fmt.Println("It's equal to 100")
	default:
		fmt.Println("It's not equal to 98 or 100")
	}
}

func sw2() {
	var num1 int = 7

	switch {
	case num1 < 0:
		fmt.Println("Number is negative")
	case num1 > 0 && num1 < 10:
		fmt.Println("Number is between 0 and 10")
	default:
		fmt.Println("Number is 10 or greater")
	}
}

func sw3() {
	k := 6
	switch k {
	case 4:
		fmt.Println("was <= 4"); fallthrough
	case 5:
		fmt.Println("was <= 5"); fallthrough
	case 6:
		fmt.Println("was <= 6"); fallthrough
	case 7:
		fmt.Println("was <= 7"); fallthrough
	case 8:
		fmt.Println("was <= 8"); fallthrough
	default:
		fmt.Println("default case")
	}
}
