package grpc

import (
	"github.com/go-playground/validator/v10"
	"github.com/si3nloong/webhook/grpc/proto"
	"github.com/si3nloong/webhook/pubsub"
	"github.com/si3nloong/webhook/shared"
)

type Server struct {
	shared.Server
	proto.UnimplementedCurlHookServiceServer
}

func NewServer(mq pubsub.MessageQueue, v *validator.Validate) proto.CurlHookServiceServer {
	svr := new(Server)
	svr.Validate = v
	svr.MessageQueue = mq
	return svr
}
