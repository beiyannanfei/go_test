package main

import (
	"fmt"
	"log"
	"bufio"
	"os"
	"github.com/nsqio/go-nsq"
)

func main() {
	strIP1 := "127.0.0.1:4150"
	strIP2 := "127.0.0.1:4152"

	producer1, err := initProducer(strIP1)
	if err != nil {
		log.Fatal("init producer1 error:", err)
	}
	producer2, err := initProducer(strIP2)
	if err != nil {
		log.Fatal("init producer2 error:", err)
	}

	defer producer1.Stop()
	defer producer2.Stop()

	//读取控制台输入
	reader := bufio.NewReader(os.Stdin)

	count := 0
	for {
		fmt.Print("please say:")
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			fmt.Println("stop producer!")
			return
		}
		if count%2 == 0 {
			err := producer1.public("test1", command)
			if err != nil {
				log.Fatal("producer1 public error:", err)
			}
		} else {
			err := producer2.public("test2", command)
			if err != nil {
				log.Fatal("producer2 public error:", err)
			}
		}

		count++
	}
}

type nsqProducer struct {
	*nsq.Producer
}

//初始化生产者
func initProducer(addr string) (*nsqProducer, error) {
	fmt.Println("init producer address:", addr)
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return &nsqProducer{producer}, nil
}

//发布消息
func (np *nsqProducer) public(topic, message string) error {
	err := np.Publish(topic, []byte(message))
	if err != nil {
		log.Println("nsq public error:", err)
		return err
	}
	return nil
}
