package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"

	"github.com/si3nloong/webhook/app/http/graphql/graph/generated"
	"github.com/si3nloong/webhook/app/http/graphql/graph/model"
	"github.com/si3nloong/webhook/app/http/graphql/transformer"
)

func (r *queryResolver) Webhooks(ctx context.Context, after *string, before *string, first *uint, last *uint, filter json.RawMessage) (*model.WebhookConnection, error) {
	limit := uint(100)
	if first != nil {
		limit = *first
	} else if last != nil {
		limit = *last
	}

	datas, _, err := r.GetWebhooks(ctx, "", limit)
	if err != nil {
		return nil, err
	}

	return transformer.ToWebhookConnection(datas), nil
}

func (r *queryResolver) Webhook(ctx context.Context, id string) (*model.Webhook, error) {
	data, err := r.FindWebhook(ctx, id)
	if err != nil {
		return nil, err
	}

	return transformer.ToWebhook(data), nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
