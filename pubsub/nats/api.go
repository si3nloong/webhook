package nats

import (
	"context"

	pb "github.com/si3nloong/signaller/grpc/proto"
	"google.golang.org/protobuf/proto"
)

func (c *Client) Publish(ctx context.Context, req *pb.SendWebhookRequest) error {
	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	if _, err := c.js.Publish(c.subj, b); err != nil {
		return err
	}

	return nil
}
