package nats

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Client struct {
}

func (c *Client) Publish(ctx context.Context) error {
	log.Println(ctx)
	return nil
}

func New() *Client {

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	// Create JetStream Context
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
	streamName := "webhook"
	stream, err := js.StreamInfo(streamName)
	if err == nats.ErrStreamNotFound {
		if _, err := js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"webhook"},
		}); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	log.Println(stream, err)
	js.QueueSubscribe(
		"webhook",
		"webhook",
		func(msg *nats.Msg) {
			// acknowledge the nats
			defer msg.Ack()
		},
		nats.AckWait(5*time.Second),
		// nats.DurableName("webhook"),
		nats.MaxDeliver(5),
	)
	// // defer nc.Close()
	// // Connect to a server
	// sc, err := stan.Connect("test", "xxx", stan.NatsConn(nc))
	// if err != nil {
	// 	panic(err)
	// }

	// stan.AckWait(20 * time.Second)
	// sc.QueueSubscribe(
	// 	"webhook",
	// 	"webhook",
	// 	func(msg *stan.Msg) {},
	// 	stan.AckWait(5*time.Second),
	// 	stan.DurableName("webhook"),
	// 	stan.MaxInflight(5),
	// )
	return &Client{}
}
