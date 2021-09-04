package pubsub

import "context"

type MessageQueue interface {
	Publish(context.Context) error
}
