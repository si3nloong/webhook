package db

import (
	"context"

	"github.com/si3nloong/webhook/model"
)

type Repository interface {
	InsertLog(ctx context.Context, data *model.Log) error
	GetLogs(ctx context.Context) ([]model.Log, error)
	FindLog(ctx context.Context, id string) (*model.Log, error)

	// GetStats(ctx context.Context) error
	// IncrStat(ctx context.Context) error
}
