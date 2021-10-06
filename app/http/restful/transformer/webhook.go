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
	noOfRetries := len(data.Retries)
	if noOfRetries > 0 {
		o.LastStatusCode = data.Retries[noOfRetries-1].Response.StatusCode
	}
	o.CreatedAt = dto.DateTime(data.CreatedAt)
	o.UpdatedAt = dto.DateTime(data.UpdatedAt)
	return
}

func ToWebhookDetail(data *entity.WebhookRequest) (o *dto.WebhookDetail) {
	o = new(dto.WebhookDetail)
	o.Webhook = *ToWebhook(data)
	o.NoOfRetries = len(data.Retries)
	o.Retries = make([]dto.WebhookRetry, 0)
	for _, r := range data.Retries {
		o.Retries = append(o.Retries, dto.WebhookRetry{
			Body:       r.Response.Body,
			StatusCode: r.Response.StatusCode,
			CreatedAt:  dto.DateTime(r.CreatedAt),
		})
	}
	return
}
