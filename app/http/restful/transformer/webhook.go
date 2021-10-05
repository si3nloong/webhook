package transformer

import (
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/app/http/restful/dto"
)

func ToWebhook(data *entity.WebhookRequest) (o *dto.Webhook) {
	o = new(dto.Webhook)
	o.ID = data.ID.String()
	o.URL = data.URL
	o.Method = data.Method
	o.Headers = make(map[string]string)
	o.Body = data.Body
	o.LastStatusCode = data.LastStatusCode
	o.CreatedAt = dto.DateTime(data.CreatedAt)
	o.UpdatedAt = dto.DateTime(data.UpdatedAt)
	return
}

func ToWebhookDetail(data *entity.WebhookRequest, retries []entity.Retry) (o *dto.WebhookDetail) {
	o = new(dto.WebhookDetail)
	o.Webhook = *ToWebhook(data)
	o.NoOfRetries = len(retries)
	o.Retries = make([]interface{}, 0)
	return
}
