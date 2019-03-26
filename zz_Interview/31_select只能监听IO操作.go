package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	ch <- 1
	select {
	case msg := <-ch:
		fmt.Println(msg)
	default://todo default可以判断chan是否已经满了 因为ch中没有写入数据，为空，所以 case不会读取成功。 则 select 执行 default 语句。
		fmt.Println("default")
	}
	//todo select 的代码形式和 switch 非常相似， 不过 select 的 case 里的操作语句只能是”IO操作”（不仅仅是取值<-channel，赋值channel<-也可以）， select 会一直等待等到某个 case 语句完成，也就是等到成功从channel中读到数据。 则 select 语句结束
	fmt.Println("success")
}
