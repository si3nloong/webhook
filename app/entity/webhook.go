package entity

import (
	"time"

	"github.com/segmentio/ksuid"
)

type WebhookRequestStatus int

const (
	WebhookRequestPending WebhookRequestStatus = iota
	WebhookRequestFailed
	WebhookRequestSuccess
)

// WebhookRequest :
type WebhookRequest struct {
	ID        ksuid.KSUID       `json:"id"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Timeout   uint              `json:"timeout"`
	Attempts  []Retry           `json:"attempts"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}
