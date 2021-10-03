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
	// Timeout    time.Duration        `json:"duration"`
	// Concurrent uint                 `json:"concurrent"`
	// Retries   uint                 `json:"retries"`
	ID        ksuid.KSUID          `json:"id"`
	URL       string               `json:"url"`
	Method    string               `json:"method"`
	Headers   map[string]string    `json:"headers"`
	Body      string               `json:"body"`
	Status    WebhookRequestStatus `json:"status"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
}
