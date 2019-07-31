package main

import (
	"github.com/nsqio/go-nsq"
	"time"
	"fmt"
)

//https://godoc.org/github.com/nsqio/go-nsq#Producer

func main() {
	nsqIp := "127.0.0.1:4150"
	producer, _ := nsq.NewProducer(nsqIp, nsq.NewConfig())
	defer producer.Stop()

	fmt.Printf("%#v", producer)

	producer.DeferredPublish("test", 10*time.Second, []byte("test-message-DeferredPublish"))
	time.Sleep(time.Second)
}
