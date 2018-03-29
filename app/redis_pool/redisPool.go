package redis_pool

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"reflect"
)

func init() {
	/**
	redis 配置
	 */
	host := "127.0.0.1"
	port := 6379
	db := 6
	pwd := ""

	redisAddr := fmt.Sprintf("%s:%d", host, port)
	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}

		if db > 0 { //选库
			if err := client.Cmd("SELECT", db).Err; err != nil {
				client.Close() //important
				return nil, err
			}
		}

		if pwd != "" { //设置密码
			if err := client.Cmd("AUTH", pwd).Err; err != nil {
				client.Close() //important
				return nil, err
			}
		}

		return client, nil
	}

	rcPool, err := pool.NewCustom("tcp", redisAddr, 10, df) //自定义redis链接库
	if err != nil {
		fmt.Printf("pool NewCustom err: %v\n", err.Error())
		return
	}

	redisClient, err := rcPool.Get()	//从连接池获取连接
	if err != nil {
		fmt.Printf("rcPool get err: %v", err.Error())
		return
	}

	setResponse, err := redisClient.Cmd("SET", "AA", "1000").Str()
	if err != nil {
		fmt.Printf("setResponse err: %v", err.Error())
		return
	}
	fmt.Printf("setResponse type: %v, value: %v\n", reflect.TypeOf(setResponse), setResponse)

	rcPool.Put(redisClient)	//将连接放回线程池   一般用法：defer rcPool.Put(redisClient)
}
