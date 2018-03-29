package redis

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"reflect"
)

var rediClient *redis.Client

func init() {
	host := "127.0.0.1"
	port := 6379

	redisAddr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("client init redis adrr: %v\n", redisAddr)
	fmt.Printf("redisClient's type is: %v\n", reflect.TypeOf(rediClient))
}
