package main

import (
	"github.com/monnand/goredis"
	"fmt"
)

func main() {
	var client0 goredis.Client
	keyStr := "go_test"
	err := client0.Set(keyStr, []byte("hello"))
	if err != nil {
		fmt.Println("redis set error:", err)
	}

	val, err := client0.Get(keyStr)
	if err != nil {
		fmt.Println("redis get error:", err)
	}
	fmt.Println("redis get val:", string(val))

	fmt.Println("---------------- select db ----------------")
	var client1 goredis.Client
	client1.Addr = "127.0.0.1:6379"
	client1.Db = 1

	keyStr = "go_l"
	vals := []string{"a", "b", "c", "d"}
	for _, v := range vals {
		client1.Rpush(keyStr, []byte(v))
	}

	dbvals, _ := client1.Lrange(keyStr, 0, -1)
	for i, v := range dbvals {
		fmt.Printf("%v = %v\n", i, string(v))
	}

	res, _ := client1.Del(keyStr)
	fmt.Println("res:", res)
}
