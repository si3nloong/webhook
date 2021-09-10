package redis

import (
	"context"
	"fmt"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"
	"github.com/si3nloong/webhook/cmd"
	pb "github.com/si3nloong/webhook/grpc/proto"
)

type messageQueue struct {
	err    chan error
	conn   rmq.Connection
	client redis.Cmdable
}

func New(ctx context.Context, cfg *cmd.Config) *messageQueue {
	q := new(messageQueue)

	if cfg.MessageQueue.Redis.Cluster {
		q.client = redis.NewClusterClient(&redis.ClusterOptions{
			Password: cfg.MessageQueue.Redis.Password,
		})
	} else {
		q.client = redis.NewClient(&redis.Options{
			Addr:     cfg.MessageQueue.Redis.Addr,
			Password: cfg.MessageQueue.Redis.Password,
			DB:       0,
		})
	}

	if err := q.client.Ping(ctx); err != nil {
		panic(err)
	}

	// c := new(redis.ClusterClient)
	// rmq.OpenConnectionWithRmqRedisClient("", c, q.err)
	conn, err := rmq.OpenConnection("my service", "tcp", "localhost:6379", 1, q.err)
	if err != nil {
		panic(err)
	}

	q.conn = conn
	return q
}

func (q *messageQueue) Publish(ctx context.Context, req *pb.SendWebhookRequest) error {
	return nil
}

func ExampleClient() {
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
