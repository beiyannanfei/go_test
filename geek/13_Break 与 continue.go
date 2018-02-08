package main

func main() {
	break1()
	println()
	continue1()

}

func continue1() {
	for i := 0; i < 10; i++ {
		if i == 5 {
			continue	//关键字 continue 忽略剩余的循环体而直接进入下一次循环的过程
		}
		print(i)
		print(" ")
	}
}

func break1() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			if j > 5 {
				break	//break 只会退出最内层的循环
			}
			print(j)
		}
		print("  ")
	}
}
