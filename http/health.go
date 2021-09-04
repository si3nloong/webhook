package http

import (
	"context"

	"github.com/si3nloong/curlhook/grpc/proto"
)

func (s *Server) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	resp := new(proto.HealthCheckResponse)
	return resp, nil
}
