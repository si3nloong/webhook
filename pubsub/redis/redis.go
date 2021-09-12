package redis

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"github.com/si3nloong/webhook/cmd"
	pb "github.com/si3nloong/webhook/grpc/proto"
	"google.golang.org/protobuf/proto"
)

type messageQueue struct {
	err    chan error
	client redis.Cmdable
	conn   rmq.Connection
	queue  rmq.Queue
}

func New(cfg cmd.Config) *messageQueue {
	mq := new(messageQueue)
	mq.err = make(chan error)

	if cfg.MessageQueue.Redis.Cluster {
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(cfg.MessageQueue.Redis.Addr, ","),
			Username: cfg.MessageQueue.Redis.Username,
			Password: cfg.MessageQueue.Redis.Password,
		})
		mq.client = clusterClient
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.MessageQueue.Redis.Addr,
			Username: cfg.MessageQueue.Redis.Username,
			Password: cfg.MessageQueue.Redis.Password,
			DB:       cfg.MessageQueue.Redis.DB,
		})
		mq.client = client
	}

	if err := mq.client.Ping(context.TODO()).Err(); err != nil {
		panic(err)
	}

	conn, err := rmq.OpenConnectionWithRedisClient(cfg.MessageQueue.Topic, mq.client.(*redis.Client), mq.err)
	if err != nil {
		panic(err)
	}

	q, err := conn.OpenQueue(cfg.MessageQueue.QueueGroup)
	if err != nil {
		panic(err)
	}

	if err := q.StartConsuming(3, 3*time.Second); err != nil {
		panic(err)
	}

	for i := 0; i < cfg.NoOfWorker; i++ {
		name, err := q.AddConsumerFunc(cfg.MessageQueue.QueueGroup, func(d rmq.Delivery) {
			log.Println("Consumer ", i)
			log.Println(d.Payload())

			req := new(pb.SendWebhookRequest)
			if err := proto.Unmarshal([]byte(d.Payload()), req); err != nil {
				return
			}
		})
		log.Println(name, err)
	}

	mq.conn = conn
	mq.queue = q
	return mq
}

func (mq *messageQueue) Publish(ctx context.Context, req *pb.SendWebhookRequest) error {
	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	if err := mq.queue.PublishBytes(b); err != nil {
		return err
	}

	return nil
}

func (mq *messageQueue) GracefulStop() error {
	mq.queue.Destroy()
	switch mq.queue.StopConsuming() {

	}
	return nil
}
