package main

import (
	"log"
	"strconv"
	"time"

	"github.com/go-micro/plugins/v4/broker/rabbitmq"
	broker "go-micro.dev/v4/broker"
)

func main() {
	rabbitmqUrl := "amqp://guest:guest@127.0.0.1:5672/"
	exchangeName := "amq.topic"
	routeTopic := "test"

	// 创建 RabbitMQ Broker
	b := rabbitmq.NewBroker(
		broker.Addrs(rabbitmqUrl),           // RabbitMQ访问地址，含VHost
		rabbitmq.ExchangeName(exchangeName), // 交换机的名称
		rabbitmq.DurableExchange(),          // 消息在Exchange中时会进行持久化处理
	)

	// 先连接，才能用
	err := b.Connect()
	if err != nil {
		log.Println(err)
	}

	loopPublish(b, routeTopic)

	log.Println("Publisher is working ...")

	select {}
}

func loopPublish(b broker.Broker, topic string) {
	tick := time.NewTicker(time.Second)
	for range tick.C {
		publish(b, topic)
	}
}

func publish(b broker.Broker, topic string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("error:", err)
		}
	}()

	log.Println("b.Publish")

	curUnix := strconv.FormatInt(time.Now().Unix(), 10)
	omsg := "{\"Id\":" + curUnix + ",\"Name\":\"张三\"}"
	msg := &broker.Message{
		Body: []byte(omsg),
	}

	err := b.Publish(topic, msg)
	if err != nil {
		log.Println(err)
	}
}
