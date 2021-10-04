package grpc

import (
	"github.com/si3nloong/webhook/app/entity"
	pb "github.com/si3nloong/webhook/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toWebhookProto(wh *entity.WebhookRequest) (proto *pb.Webhook) {
	proto = new(pb.Webhook)
	proto.Id = wh.ID.String()
	proto.Body = wh.Body
	proto.Method = wh.Method
	proto.Retries = uint32(len(wh.Retries))
	proto.CreatedAt = timestamppb.New(wh.CreatedAt)
	proto.UpdatedAt = timestamppb.New(wh.UpdatedAt)
	return
}
