package entity

import (
	"time"

	"github.com/segmentio/ksuid"
)

// Log :
type Log struct {
	ID        ksuid.KSUID       `json:"id"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Method    string            `json:"method"`
	CreatedAt time.Time         `json:"createdAt"`
}
