package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/si3nloong/webhook/app/shared"
	"github.com/si3nloong/webhook/pubsub"
)

type Server struct {
	shared.Server
}

func NewServer(mq pubsub.MessageQueue, v *validator.Validate) *Server {
	svr := new(Server)
	svr.Validate = v
	svr.MessageQueue = mq
	return svr
}
