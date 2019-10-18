package main

import (
	"sync"
	"fmt"
)

type UserAges struct {
	ages map[string]int
	sync.Mutex	//互斥锁
}

func (ua *UserAges) Add(name string, age int) {
	ua.Lock()
	defer ua.Unlock()
	ua.ages[name] = age
}

func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok {
		return age
	}

	return -1
}

func main() {
	count := 1000
	gw := sync.WaitGroup{}
	gw.Add(count * 3)
	u := UserAges{ages: map[string]int{}}
	add := func(i int) {
		u.Add(fmt.Sprintf("user_%d", i), i)
		gw.Done()
	}

	for i := 0; i < count; i++ {
		go add(i)
		go add(i)
	}

	for i := 0; i < count; i++ {
		go func(i int) {
			defer gw.Done()
			u.Get(fmt.Sprintf("user_%d", i))
		}(i)
	}

	gw.Wait()
	fmt.Println("Done")
}

//虽然有使用sync.Mutex做写锁，但是map是并发读写不安全的。map属于引用类型，并发读写时多个协程见是通过指针访问同一个地址，即访问共享变量，
// 此时同时读写资源存在竞争关系。所以会报错误信息:fatal error: concurrent map read and map write。
