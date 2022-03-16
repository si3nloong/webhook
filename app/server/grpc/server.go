package grpc

import (
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/si3nloong/webhook/app/shared"
	"github.com/si3nloong/webhook/cmd"
	pb "github.com/si3nloong/webhook/protobuf"
	"google.golang.org/grpc"
)

type Server struct {
	ws shared.WebhookServer
	pb.UnimplementedWebhookServiceServer
}

func NewServer(cfg *cmd.Config, ws shared.WebhookServer) *grpc.Server {
	opts := make([]grpc.ServerOption, 0)
	if cfg.GRPC.ApiKey != "" {
		opts = append(opts, grpc.UnaryInterceptor(grpcauth.UnaryServerInterceptor(authorizationInterceptor(cfg.GRPC.ApiKey))))
	}

	grpcServer := grpc.NewServer(opts...)
	svr := &Server{ws: ws}
	pb.RegisterWebhookServiceServer(grpcServer, svr)
	return grpcServer
}
