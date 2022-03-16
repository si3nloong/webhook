package transformer

import (
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/app/server/rest/dto"
)

func ToWebhook(data *entity.WebhookRequest) (o *dto.Webhook) {
	o = new(dto.Webhook)
	o.ID = data.ID.String()
	o.URL = data.URL
	o.Method = data.Method
	o.Headers = make(map[string]string)
	o.Body = data.Body
	o.Timeout = data.Timeout
	noOfRetries := len(data.Attempts)
	if noOfRetries > 0 {
		o.LastStatusCode = data.Attempts[noOfRetries-1].StatusCode
	}
	o.CreatedAt = dto.DateTime(data.CreatedAt)
	o.UpdatedAt = dto.DateTime(data.UpdatedAt)
	return
}

func ToWebhookDetail(data *entity.WebhookRequest) (o *dto.WebhookDetail) {
	o = new(dto.WebhookDetail)
	o.Webhook = *ToWebhook(data)
	o.NoOfRetries = len(data.Attempts)
	o.Attempts = make([]dto.WebhookRetry, 0)
	for _, r := range data.Attempts {
		attempt := dto.WebhookRetry{}
		attempt.Headers = make(map[string]string)
		attempt.Body = r.Body
		attempt.ElapsedTime = r.ElapsedTime
		attempt.StatusCode = r.StatusCode
		attempt.CreatedAt = dto.DateTime(r.CreatedAt)

		o.Attempts = append(o.Attempts, attempt)
	}
	return
}
