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
	Attempts  []Attempt         `json:"attempts"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// Attempt :
type Attempt struct {
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	StatusCode  uint              `json:"statusCode"`
	ElapsedTime int64             `json:"elapsedTime"`
	CreatedAt   time.Time         `json:"createdAt"`
}
