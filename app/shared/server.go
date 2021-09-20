package shared

import (
	"github.com/go-playground/validator/v10"
	"github.com/si3nloong/webhook/pubsub"
)

type Server struct {
	*validator.Validate
	pubsub.MessageQueue
}
