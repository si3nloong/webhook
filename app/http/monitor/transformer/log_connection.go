package transformer

import (
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/app/http/monitor/graph/model"
)

func ToLogConnection(in []*entity.WebhookRequest) (conn *model.LogConnection) {
	conn = new(model.LogConnection)
	conn.Nodes = ToLogs(in)
	conn.PageInfo = new(model.PageInfo)
	conn.TotalCount = 100
	return
}

func ToLogs(in []*entity.WebhookRequest) (out []*model.Log) {
	out = make([]*model.Log, len(in))
	for idx := range in {
		out[idx] = ToLog(in[idx])
	}
	return
}

func ToLog(in *entity.WebhookRequest) (out *model.Log) {
	out = new(model.Log)
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
