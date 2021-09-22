package grpc

import (
	"context"

	"github.com/si3nloong/webhook/app/grpc/proto"
)

func (s *Server) SendWebhook(ctx context.Context, req *proto.SendWebhookRequest) (*proto.SendWebhookResponse, error) {
	// if err := s.StructCtx(ctx, req); err != nil {
	// 	return nil, status.Convert(err).Err()
	// }

	// push to message queue
	// if err := s.Publish(ctx, req); err != nil {
	// 	return nil, status.Convert(err).Err()
	// }

	resp := new(proto.SendWebhookResponse)
	return resp, nil
}
