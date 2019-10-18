package main

import "fmt"
import "sync"
import "time"

type ThreadSafeSet struct {
	sync.RWMutex
	s []int
}

func (set *ThreadSafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		set.RLock()
		for elem := range set.s {
			ch <- elem
			fmt.Println("get:", elem, ",")
		}
		close(ch)
		set.RUnlock()
	}()
	return ch
}

func read() {
	set := ThreadSafeSet{}
	set.s = make([]int, 3)
	ch := set.Iter()
	closed := false
	for {
		select {
		case v, ok := <-ch:
			if ok {
				fmt.Println("read:", v, ",")
			} else {
				closed = true
			}
		}
		if closed {
			fmt.Println("closed")
			break
		}
	}
	fmt.Print("Done")
}

func unRead() {
	set := ThreadSafeSet{}
	set.s = make([]int, 100)
	ch := set.Iter()
	_ = ch
	time.Sleep(5 * time.Second)
	fmt.Print("Done")
}

func main() {
	//read()
	unRead()
}

//内部迭代出现阻塞。默认初始化时无缓冲区，需要等待接收者读取后才能继续写入。
//chan在使用make初始化时可附带一个可选参数来设置缓冲区。默认无缓冲，题目中便初始化的是无缓冲区的chan，这样只有写入的元素直到被读取后才能继续写入，不然就一直阻塞。
//设置缓冲区大小后，写入数据时可连续写入到缓冲区中，直到缓冲区被占满。从chan中接收一次便可从缓冲区中释放一次。可以理解为chan是可以设置吞吐量的处理池。
