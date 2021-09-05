package pubsub

import (
	"context"

	"github.com/si3nloong/signaller/grpc/proto"
)

type MessageQueue interface {
	Publish(context.Context, *proto.SendWebhookRequest) error
}
