package pubsub

import (
	"context"

	"github.com/si3nloong/webhook/grpc/proto"
)

type ConsumerFunc func(*proto.SendWebhookRequest) error

type MessageQueue interface {
	Publish(ctx context.Context, req *proto.SendWebhookRequest) error
	// SubscribeOn(func([]byte))
	GracefulStop() error
}
