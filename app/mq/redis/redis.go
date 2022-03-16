package redis

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/cmd"
)

type redisMQ struct {
	err     chan error
	subs    []string
	client  redis.Cmdable
	conn    rmq.Connection
	queue   rmq.Queue
	cleaner *rmq.Cleaner
}

func New(cfg *cmd.Config, cb func(delivery rmq.Delivery)) (*redisMQ, error) {
	mq := new(redisMQ)
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

	conn, err := rmq.OpenConnectionWithRedisClient(cfg.MessageQueue.Topic, mq.client, mq.err)
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

	// setup consumers
	for i := 0; i < cfg.NoOfWorker; i++ {
		name, err := q.AddConsumerFunc(cfg.MessageQueue.QueueGroup, cb)
		if err != nil {
			panic(err)
		}

		mq.subs = append(mq.subs, name)
	}

	mq.cleaner = rmq.NewCleaner(conn)
	mq.conn = conn
	mq.queue = q
	return mq, nil
}

func (mq *redisMQ) Publish(ctx context.Context, data *entity.WebhookRequest) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := mq.queue.PublishBytes(b); err != nil {
		return err
	}

	return nil
}

// func (mq *redisMQ) GracefulStop() error {
// 	mq.queue.Destroy()
// 	switch mq.queue.StopConsuming() {

// 	}
// 	return nil
// }
