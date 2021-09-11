package redis

import (
	"time"

	"github.com/adjust/rmq/v4"
)

type taskConsumer struct {
	name          string
	AutoAck       bool
	AutoFinish    bool
	SleepDuration time.Duration

	LastDelivery   rmq.Delivery
	LastDeliveries []rmq.Delivery

	finish chan int
}

func NewTestConsumer(name string) rmq.Consumer {
	return &taskConsumer{
		name:       name,
		AutoAck:    true,
		AutoFinish: true,
		finish:     make(chan int),
	}
}

func (c *taskConsumer) String() string {
	return c.name
}

func (c *taskConsumer) Consume(delivery rmq.Delivery) {
	c.LastDelivery = delivery
	c.LastDeliveries = append(c.LastDeliveries, delivery)

	if c.SleepDuration > 0 {
		time.Sleep(c.SleepDuration)
	}
	if c.AutoAck {
		if err := delivery.Ack(); err != nil {
			panic(err)
		}
	}
	if !c.AutoFinish {
		<-c.finish
	}
}

func (c *taskConsumer) Finish() {
	c.finish <- 1
}

func (c *taskConsumer) FinishAll() {
	close(c.finish)
}
