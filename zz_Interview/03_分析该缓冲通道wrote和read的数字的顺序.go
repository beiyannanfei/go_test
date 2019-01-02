package main

import (
	"log"
	"time"
)

//https://studygolang.com/articles/17374

func write(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
		log.Printf("wrote value %v to ch.", i)
	}
	close(ch)
}

func main() {
	ch := make(chan int, 2)
	go write(ch)
	time.Sleep(2 * time.Second)
	for v := range ch {
		log.Printf("read value: %v from ch.", v)
		time.Sleep(2 * time.Second)
	}
}

/*
输出结果:
在main中声明的ch的缓冲容量为2，根据缓冲通道的特点，当通道写满的时候写入方法write就会进入阻塞，
range方法会读取通道ch中的值，由于存在2s的sleep，最终结果为
2019/01/02 14:18:34 wrote value 0 to ch.
2019/01/02 14:18:34 wrote value 1 to ch.
2019/01/02 14:18:36 read value: 0 from ch.
2019/01/02 14:18:36 wrote value 2 to ch.
2019/01/02 14:18:38 read value: 1 from ch.
2019/01/02 14:18:38 wrote value 3 to ch.
2019/01/02 14:18:40 read value: 2 from ch.
2019/01/02 14:18:40 wrote value 4 to ch.
2019/01/02 14:18:42 read value: 3 from ch.
2019/01/02 14:18:44 read value: 4 from ch.

*/