package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

// Log :
type Log struct {
	ID        ksuid.KSUID `json:"id"`
	URL       string      `json:"url"`
	Method    string      `json:"method"`
	CreatedAt time.Time   `json:"createdAt"`
}
