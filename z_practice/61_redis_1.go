package main

import (
	"github.com/go-redis/redis"
	"fmt"
)

func main() {
	conf := redis.Options{
		Addr: "127.0.0.1:6379",
	}
	conf.DB = 1

	client := redis.NewClient(&conf)

	res := client.Set("aaaa", 123, 0)
	fmt.Println(res.Val())
	fmt.Println(res.Result())
	fmt.Println(res.String())
	fmt.Println(res.Err())
	fmt.Println("---------------------------")

	val, err := client.Get("aaaa").Result()
	fmt.Println(val, err)

	val, err = client.Get("aaaa1").Result()
	if err == redis.Nil {
		fmt.Println("key is not exists")
	}
}
