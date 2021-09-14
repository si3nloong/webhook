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
	config       cmd.Config
	err          chan error
	subsriptions []string
	client       redis.Cmdable
	conn         rmq.Connection
	queue        rmq.Queue
	cleaner      *rmq.Cleaner
}

func New(cfg cmd.Config) (*messageQueue, error) {
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

	if err := mq.client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	conn, err := rmq.OpenConnectionWithRedisClient(cfg.MessageQueue.Topic, mq.client.(*redis.Client), mq.err)
	if err != nil {
		return nil, err
	}

	q, err := conn.OpenQueue(cfg.MessageQueue.QueueGroup)
	if err != nil {
		return nil, err
	}

	if err := q.StartConsuming(3, 3*time.Second); err != nil {
		return nil, err
	}

	mq.cleaner = rmq.NewCleaner(conn)
	mq.conn = conn
	mq.queue = q
	return mq, nil
}

func (mq *messageQueue) Publish(ctx context.Context, req *pb.SendWebhookRequest) error {
	log.Println("publishing 0")
	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	log.Println("publishing 1")

	if err := mq.queue.PublishBytes(b); err != nil {
		return err
	}

	return nil
}

func (mq *messageQueue) SubscribeOn(func()) {
	for i := 0; i < mq.config.NoOfWorker; i++ {
		name, err := mq.queue.AddConsumer(
			mq.config.MessageQueue.QueueGroup,
			newTaskConsumer(func(r *pb.SendWebhookRequest) error {
				log.Println(r)
				return nil
			}),
		)
		if err != nil {
			panic(err)
		}
		mq.subsriptions = append(mq.subsriptions, name)
	}
}

func (mq *messageQueue) GracefulStop() error {
	mq.queue.Destroy()
	switch mq.queue.StopConsuming() {

	}
	return nil
}
