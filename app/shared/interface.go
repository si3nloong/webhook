package shared

import (
	"context"

	"github.com/si3nloong/webhook/app/entity"
	pb "github.com/si3nloong/webhook/protobuf"
)

type Repository interface {
	CreateWebhook(ctx context.Context, data *entity.WebhookRequest) error
	GetWebhooks(ctx context.Context, curCursor string, limit uint) (datas []*entity.WebhookRequest, nextCursor string, totalCount int64, err error)
	FindWebhook(ctx context.Context, id string) (*entity.WebhookRequest, error)
	UpdateWebhook(ctx context.Context, id string, attempt *entity.Attempt) error
}

type MessageQueue interface {
	Publish(ctx context.Context, data *entity.WebhookRequest) error
}

type WebhookServer interface {
	Repository
	Validate(src interface{}) error
	Publish(ctx context.Context, req *pb.SendWebhookRequest) (*entity.WebhookRequest, error)
	LogError(err error)

	// TODO: better name (rename please)
	VarCtx(ctx context.Context, src interface{}, tag string) error
}
