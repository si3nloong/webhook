package grpc

import (
	"context"
	"log"

	pb "github.com/si3nloong/webhook/protobuf"
	"google.golang.org/grpc/status"
)

func (s *Server) GetWebhooks(ctx context.Context, req *pb.ListWebhooksRequest) (*pb.ListWebhooksResponse, error) {
	if err := s.ws.Validate(req); err != nil {
		return nil, status.Convert(err).Err()
	}

	datas, nextCursor, _, err := s.ws.GetWebhooks(ctx, req.PageToken, 100)
	if err != nil {
		return nil, status.Convert(err).Err()
	}

	resp := new(pb.ListWebhooksResponse)
	for _, data := range datas {
		resp.Webhooks = append(resp.Webhooks, toWebhookProto(data))
	}
	resp.NextPageToken = nextCursor
	return resp, nil
}

func (s *Server) GetWebhook(ctx context.Context, req *pb.GetWebhookRequest) (*pb.GetWebhookResponse, error) {
	if err := s.ws.Validate(req); err != nil {
		return nil, status.Convert(err).Err()
	}

	data, err := s.ws.FindWebhook(ctx, req.Id)
	if err != nil {
		return nil, status.Convert(err).Err()
	}

	resp := new(pb.GetWebhookResponse)
	resp.Webhook = toWebhookProto(data)
	return resp, nil
}

func (s *Server) SendWebhook(ctx context.Context, req *pb.SendWebhookRequest) (*pb.SendWebhookResponse, error) {
	if err := s.ws.Validate(req); err != nil {
		return nil, status.Convert(err).Err()
	}

	// push to message queue
	data, err := s.ws.Publish(ctx, req)
	if err != nil {
		return nil, status.Convert(err).Err()
	}

	log.Println(data)
	resp := new(pb.SendWebhookResponse)
	return resp, nil
}
