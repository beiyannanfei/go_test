package redis

//https://godoc.org/github.com/mediocregopher/radix.v2/redis
import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"reflect"
)

var redisClient *redis.Client

func init() {
	/**
	redis 配置
	 */
	host := "127.0.0.1"
	port := 6379
	db := 6

	redisAddr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("client init redis adrr: %v\n", redisAddr)
	fmt.Printf("redisClient's type is: %v\n", reflect.TypeOf(redisClient))
	var err error
	redisClient, err = redis.Dial("tcp", redisAddr)
	if err != nil {
		fmt.Printf("redis.Dial err: %v\n", err.Error())
		return
	}

	if db > 0 {
		if err = redisClient.Cmd("SELECT", db).Err; err != nil {
			fmt.Printf("redis exec select command err: %v, db: %v\n", err, db)
			return
		}
	}

	setResponse, err := redisClient.Cmd("SET", "aa", 100).Str()
	fmt.Printf("setResponse type: %v, value: %#v\n", reflect.TypeOf(setResponse), setResponse)
	set2response, err := redisClient.Cmd("SET", "bb", 200).Bytes()
	fmt.Printf("set2response type: %v, value: %#v, string(value): %v\n", reflect.TypeOf(set2response), set2response, string(set2response))

	set3response := redisClient.Cmd("SET", "cc", 300)
	set3Value, set3Err := set3response.Str()
	fmt.Printf("set3Value: %v, set3Err: %v\n", set3Value, set3Err)

	getResponse, err := redisClient.Cmd("GET", "aa").Str()
	fmt.Printf("getResponse type: %v, value: %#v\n", reflect.TypeOf(getResponse), getResponse)

	mgetResponse := redisClient.Cmd("MGET", "aa", "bb", "cc") //一次获取多个
	if mgetResponse.Err != nil {
		fmt.Printf("MGET err: %v\n", err)
		return
	}
	fmt.Printf("mgetResponse type: %v, value: %#v\n", reflect.TypeOf(mgetResponse), mgetResponse)

	mgetList, err := mgetResponse.List()
	fmt.Printf("mgetList type: %v, value: %v, err: %v\n", reflect.TypeOf(mgetList), mgetList, err)
	for elemKey, elemVal := range mgetList {
		fmt.Printf("elemKey: %#v, elemVal: %#v\n", elemKey, elemVal)
	}

	//Array() <=> List()
	mgetArray, err := mgetResponse.Array()
	fmt.Printf("mgetArray type: %v, value: %v, err: %v\n", reflect.TypeOf(mgetArray), mgetArray, err)
	for elemKey, elemObj := range mgetArray {
		val, _ := elemObj.Str()
		fmt.Printf("elemKey: %v, elemObj type: %v, elemObj value: %v, elemObj.Str: %v\n", elemKey, reflect.TypeOf(elemObj), elemObj, val)
	}

	err = redisClient.Close() //切记，用完redis要关闭
	if err != nil {
		fmt.Printf("redisClient Close err: %v", err.Error())
		return
	}
}
