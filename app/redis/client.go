package redis

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
)

var rediClient *redis.Client

func init() {
	host := "127.0.0.1"
	port := 6379

	redisAddr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("client init redis adrr: %v\n", redisAddr)

}
