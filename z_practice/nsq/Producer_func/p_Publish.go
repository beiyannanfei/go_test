package main

import (
	"github.com/nsqio/go-nsq"
	"time"
)

//https://godoc.org/github.com/nsqio/go-nsq#Producer

func main() {
	nsqIp := "127.0.0.1:4150"
	producer, _ := nsq.NewProducer(nsqIp, nsq.NewConfig())
	defer producer.Stop()

	producer.Publish("test", []byte("test-message"))
	time.Sleep(time.Second)
}
