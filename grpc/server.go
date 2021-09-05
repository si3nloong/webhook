package grpc

import (
	"github.com/go-playground/validator/v10"
	"github.com/si3nloong/signaller/grpc/proto"
	"github.com/si3nloong/signaller/pubsub"
)

type Server struct {
	*validator.Validate
	mq pubsub.MessageQueue
	proto.UnimplementedCurlHookServiceServer
}

func NewServer(v *validator.Validate, mq pubsub.MessageQueue) proto.CurlHookServiceServer {
	return &Server{
		Validate: v,
		mq:       mq,
	}
}
