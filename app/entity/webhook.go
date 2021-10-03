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
	ID        ksuid.KSUID          `json:"id"`
	URL       string               `json:"url"`
	Method    string               `json:"method"`
	Headers   map[string]string    `json:"headers"`
	Body      string               `json:"body"`
	Retries   uint                 `json:"retries"`
	Status    WebhookRequestStatus `json:"status"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	// Timeout    time.Duration        `json:"duration"`
	// Concurrent uint                 `json:"concurrent"`
}
