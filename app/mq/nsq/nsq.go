package nsq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type nsqMQ struct {
}

func New() *nsqMQ {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println(consumer)
	return &nsqMQ{}
}
