package grpc

import (
	"context"

	pb "github.com/si3nloong/webhook/protobuf"
)

func (s *Server) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	resp := new(pb.HealthCheckResponse)
	resp.Status = pb.HealthCheckResponse_SERVING
	return resp, nil
}

func (s *Server) Watch(req *pb.HealthCheckRequest, stream pb.WebhookService_WatchServer) error {
	resp := new(pb.HealthCheckResponse)
	resp.Status = pb.HealthCheckResponse_SERVING
	return stream.SendMsg(resp)
}
