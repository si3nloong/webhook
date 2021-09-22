package transformer

import (
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/app/http/monitor/graph/model"
)

func ToLogConnection(in []entity.Log) (conn *model.LogConnection) {
	conn = new(model.LogConnection)
	conn.Nodes = ToLogs(in)
	conn.TotalCount = 100
	return
}

func ToLogs(in []entity.Log) (out []*model.Log) {
	out = make([]*model.Log, len(in))
	for idx := range in {
		out[idx] = ToLog(in[idx])
	}
	return
}

func ToLog(in entity.Log) (out *model.Log) {
	out = new(model.Log)
	out.ID = in.ID.String()
	out.URL = in.URL
	out.Method = model.HTTPMethod(in.Method)
	out.Body = in.Body
	out.CreatedAt = in.CreatedAt
	return
}
