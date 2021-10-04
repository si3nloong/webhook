package redis

import (
	"log"
	"time"

	"github.com/adjust/rmq/v4"
	pb "github.com/si3nloong/webhook/protobuf"
	"google.golang.org/protobuf/proto"
)

type taskConsumer struct {
	name          string
	AutoAck       bool
	AutoFinish    bool
	SleepDuration time.Duration
	cb            func(*pb.SendWebhookRequest) error

	LastDelivery   rmq.Delivery
	LastDeliveries []rmq.Delivery

	finish chan int
}

func newTaskConsumer(cb func(*pb.SendWebhookRequest) error) rmq.Consumer {
	return &taskConsumer{
		// name:       name,
		cb:         cb,
		AutoAck:    true,
		AutoFinish: true,
		finish:     make(chan int),
	}
}

func (c *taskConsumer) String() string {
	return c.name
}

func (c *taskConsumer) Consume(delivery rmq.Delivery) {
	req := new(pb.SendWebhookRequest)
	if err := proto.Unmarshal([]byte(delivery.Payload()), req); err != nil {
		return
	}

	log.Println("hERE 1")
	log.Println(delivery)

	if err := c.cb(req); err != nil {
		return
	}

	delivery.Ack()
	// log.Println(delivery.Reject())
	// c.LastDelivery = delivery
	// c.LastDeliveries = append(c.LastDeliveries, delivery)

	// if c.SleepDuration > 0 {
	// 	time.Sleep(c.SleepDuration)
	// }
	// if c.AutoAck {
	// 	if err := delivery.Ack(); err != nil {
	// 		panic(err)
	// 	}
	// }
	// if !c.AutoFinish {
	// 	<-c.finish
	// }
}

func (c *taskConsumer) Finish() {
	c.finish <- 1
}

func (c *taskConsumer) FinishAll() {
	close(c.finish)
}
