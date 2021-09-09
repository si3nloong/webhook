package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	pb "github.com/si3nloong/webhook/grpc/proto"
)

type messageQueue struct {
	redis redis.Cmdable
}

func (m *messageQueue) Publish(ctx context.Context, req *pb.SendWebhookRequest) error {
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
