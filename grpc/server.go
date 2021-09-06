package grpc

import (
	"github.com/go-playground/validator/v10"
	"github.com/si3nloong/rwhook/grpc/proto"
	"github.com/si3nloong/rwhook/pubsub"
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
