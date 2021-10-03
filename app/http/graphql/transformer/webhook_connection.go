package transformer

import (
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/app/http/graphql/graph/model"
)

func ToWebhookConnection(in []*entity.WebhookRequest) (conn *model.WebhookConnection) {
	conn = new(model.WebhookConnection)
	conn.Nodes = ToLogs(in)
	conn.PageInfo = new(model.PageInfo)
	conn.TotalCount = 100
	return
}

func ToLogs(in []*entity.WebhookRequest) (out []*model.Webhook) {
	out = make([]*model.Webhook, len(in))
	for idx := range in {
		out[idx] = ToWebhook(in[idx])
	}
	return
}

func ToWebhook(in *entity.WebhookRequest) (out *model.Webhook) {
	out = new(model.Webhook)
	out.ID = in.ID.String()
	out.URL = in.URL
	out.Method = model.HTTPMethod(in.Method)
	out.Headers = make([]*model.HTTPHeader, 0)
	for k, v := range in.Headers {
		out.Headers = append(out.Headers, &model.HTTPHeader{Key: k, Value: v})
	}
	out.Body = in.Body
	// out.Retries =
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	return
}
