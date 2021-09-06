package pubsub

import (
	"context"

	"github.com/si3nloong/rwhook/grpc/proto"
)

type MessageQueue interface {
	Publish(context.Context, *proto.SendWebhookRequest) error
}
