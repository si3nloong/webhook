package grpc

import (
	"context"

	"github.com/si3nloong/webhook/grpc/proto"
	"google.golang.org/grpc/status"
)

func (s *Server) SendWebhook(ctx context.Context, req *proto.SendWebhookRequest) (*proto.SendWebhookResponse, error) {
	if err := s.StructCtx(ctx, req); err != nil {
		return nil, status.Convert(err).Err()
	}

	// push to nats
	if err := s.Publish(ctx, req); err != nil {
		return nil, status.Convert(err).Err()
	}

	resp := new(proto.SendWebhookResponse)
	return resp, nil
}
