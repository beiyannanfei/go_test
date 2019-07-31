package main

import (
	"time"
	"github.com/nsqio/go-nsq"
	"github.com/beiyannanfei/go_test/z_practice/nsq"
	"fmt"
)

type nsqHandler struct {
	nsqConsumer      *nsq.Consumer
	messagesReceived int
}

//处理消息
func (nh *nsqHandler) HandleMessage(msg *nsq.Message) error {
	nh.messagesReceived++
	fmt.Printf("receive ID: %s, addr: %s, message: %s\n", msg.ID, msg.NSQDAddress, string(msg.Body))
	return nil
}

func main() {
	go nsq_1.PublishMsg("test", "test-msg111", 1, true)

	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 30 * time.Second
	c, err := nsq.NewConsumer("test", "test-channel", cfg)
	if err != nil {
		fmt.Println("NewConsumer failed: ", err.Error())
		return
	}
	defer c.Stop()

	handler := &nsqHandler{nsqConsumer: c}
	c.AddHandler(handler)

	c.ConnectToNSQLookupd("127.0.0.1:4161")

	time.Sleep(10 * time.Second)

	stats := c.Stats()
	fmt.Println(stats)
}
