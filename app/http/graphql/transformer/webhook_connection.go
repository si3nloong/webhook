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
	out.Timeout = in.Timeout
	out.Attempts = make([]*model.WebhookAttempt, len(in.Attempts))
	for idx, a := range in.Attempts {
		attempt := model.WebhookAttempt{}
		attempt.ElapsedTime = a.ElapsedTime
		attempt.Headers = make([]*model.HTTPHeader, 0)
		for k, v := range a.Headers {
			attempt.Headers = append(attempt.Headers, &model.HTTPHeader{Key: k, Value: v})
		}
		attempt.Body = a.Body
		attempt.StatusCode = uint(a.StatusCode)
		attempt.CreatedAt = a.CreatedAt

		out.Attempts[idx] = &attempt
	}
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	return
}
