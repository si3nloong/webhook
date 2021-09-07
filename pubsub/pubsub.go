package pubsub

import (
	"context"

	"github.com/si3nloong/webhook/grpc/proto"
)

type MessageQueue interface {
	Publish(context.Context, *proto.SendWebhookRequest) error
}
