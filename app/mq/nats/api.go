package nats

import (
	"context"
	"encoding/json"

	"github.com/si3nloong/webhook/app/entity"
)

func (q *natsMQ) Publish(ctx context.Context, data *entity.WebhookRequest) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := q.js.Publish(q.subj, b); err != nil {
		return err
	}
	return nil
}
