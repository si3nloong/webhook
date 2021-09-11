package grpc

import (
	"github.com/go-playground/validator/v10"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/si3nloong/webhook/cmd"
	"github.com/si3nloong/webhook/grpc/proto"
	"github.com/si3nloong/webhook/pubsub"
	"github.com/si3nloong/webhook/shared"
	"google.golang.org/grpc"
)

type Server struct {
	shared.Server
	proto.UnimplementedCurlHookServiceServer
}

func NewServer(cfg cmd.Config, mq pubsub.MessageQueue, v *validator.Validate) *grpc.Server {
	opts := make([]grpc.ServerOption, 0)
	if cfg.GRPC.ApiKey != "" {
		opts = append(opts, grpc.UnaryInterceptor(grpcauth.UnaryServerInterceptor(authorizationInterceptor(cfg.GRPC.ApiKey))))
	}

	grpcServer := grpc.NewServer(opts...)
	svr := new(Server)
	svr.Validate = v
	svr.MessageQueue = mq
	proto.RegisterCurlHookServiceServer(grpcServer, svr)
	return grpcServer
}
