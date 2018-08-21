package main

import (
	"time"
	"fmt"
)
//chan关闭的情况取超出值为类型默认值，如int为0 interface为nil
func main() {
	intChan := make(chan int, 10)

	for i := 0; i < 10; i++ {
		intChan <- i
	}

	close(intChan)

	for {
		i := <-intChan
		//十次后i值都为0，不报错
		time.Sleep(time.Second)
		fmt.Println(i)
	}
}
