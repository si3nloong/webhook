package grpc

import (
	"context"

	"github.com/si3nloong/webhook/app/grpc/proto"
)

func (s *Server) Check(ctx context.Context, req *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	resp := new(proto.HealthCheckResponse)
	resp.Status = proto.HealthCheckResponse_SERVING
	return resp, nil
}

func (s *Server) Watch(req *proto.HealthCheckRequest, stream proto.CurlHookService_WatchServer) error {
	resp := new(proto.HealthCheckResponse)
	resp.Status = proto.HealthCheckResponse_SERVING
	return stream.SendMsg(resp)
}
