package nsq_1

import (
	"github.com/nsqio/go-nsq"
	"time"
	"fmt"
)

func PublishMsg(topic string, msg string, inteval int64, isDefer bool) {
	nsqIp := "127.0.0.1:4150"
	producer, _ := nsq.NewProducer(nsqIp, nsq.NewConfig())
	defer producer.Stop()

	for {
		if isDefer {
			producer.DeferredPublish(topic, 3*time.Second, []byte(fmt.Sprintf("%s-%d", msg, time.Now().Unix())))
		} else {
			producer.Publish(topic, []byte(fmt.Sprintf("%s-%d", msg, time.Now().Unix())))
		}

		time.Sleep(time.Duration(inteval) * time.Second)
	}
}
