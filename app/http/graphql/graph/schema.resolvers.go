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
	var (
		limit     = defaultLimit
		curCursor string
	)
	if first != nil && *first <= defaultLimit {
		limit = *first
	} else if last != nil && *last <= defaultLimit {
		limit = *last
	}
	if after != nil {
		curCursor = *after
	}

	datas, nextCursor, totalCount, err := r.GetWebhooks(ctx, curCursor, limit)
	if err != nil {
		return nil, err
	}

	return transformer.ToWebhookConnection(
		datas,
		curCursor,
		nextCursor,
		totalCount,
	), nil
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const defaultLimit = uint(100)
