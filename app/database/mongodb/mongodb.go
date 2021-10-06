package mongodb

import (
	"time"
)

type db struct {
	indexName string
	timeout   time.Duration
}
