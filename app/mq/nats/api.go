package nats

import (
	"context"

	pb "github.com/si3nloong/webhook/app/grpc/proto"
	"google.golang.org/protobuf/proto"
)

func (q *natsMQ) Publish(ctx context.Context, req *pb.SendWebhookRequest) error {
	b, err := proto.Marshal(req)
	if err != nil {
		return err
	}

	if _, err := q.js.Publish(q.subj, b); err != nil {
		return err
	}

	return nil
}
