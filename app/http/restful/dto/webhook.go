package dto

// Webhook :
type Webhook struct {
	ID             string            `json:"id"`
	URL            string            `json:"url"`
	Method         string            `json:"method"`
	Body           string            `json:"body"`
	Headers        map[string]string `json:"headers"`
	Timeout        uint              `json:"timeout"`
	LastStatusCode int               `json:"lastStatusCode"`
	Status         string            `json:"status"`
	CreatedAt      DateTime          `json:"createdAt"`
	UpdatedAt      DateTime          `json:"updatedAt"`
}

// WebhookDetail :
type WebhookDetail struct {
	Webhook
	NoOfRetries int            `json:"noOfRetries"`
	Retries     []WebhookRetry `json:"retries"`
}

// WebhookRetry :
type WebhookRetry struct {
	Body       string   `json:"body"`
	StatusCode int      `json:"statusCode"`
	CreatedAt  DateTime `json:"created"`
}
