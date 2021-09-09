package pubsub

import (
	"context"

	"github.com/si3nloong/webhook/grpc/proto"
)

type MessageQueue interface {
	// Subscribe(worker uint, cb func())
	Publish(ctx context.Context, req *proto.SendWebhookRequest) error
	GracefulStop() error
}
