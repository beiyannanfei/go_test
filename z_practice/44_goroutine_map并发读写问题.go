package main

import (
	"sync"
	"time"
	"fmt"
	"strconv"
)

//三、解决方法
//１、加锁
//(1)通用锁

type Demo struct {
	Data map[string]string
	Lock sync.Mutex
}

func (d *Demo) Get(key string) string {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	return d.Data[key]
}

func (d *Demo) Set(key, value string) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Data[key] = value
}

func main() {
	v := Demo{Data: make(map[string]string), Lock: sync.Mutex{}}

	//v.Set("aa", "11")
	//fmt.Println(v.Get("aa"))

	go func() { //开一个协程写map
		for i := 0; i < 100000; i++ {
			v.Set(strconv.Itoa(i), strconv.Itoa(i))
		}
	}()

	time.Sleep(time.Second * 1)
	go func() { //开一个协程读map
		for i := 0; i < 100000; i++ {
			fmt.Println(v.Get(strconv.Itoa(i)))
		}
	}()

	time.Sleep(time.Second * 20)
}
