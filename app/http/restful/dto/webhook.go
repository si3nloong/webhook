package dto

// Webhook :
type Webhook struct {
	ID             string            `json:"id"`
	URL            string            `json:"url"`
	Method         string            `json:"method"`
	Headers        map[string]string `json:"headers"`
	Body           string            `json:"body"`
	Timeout        uint              `json:"timeout"`
	LastStatusCode int               `json:"lastStatusCode"`
	CreatedAt      DateTime          `json:"createdAt"`
	UpdatedAt      DateTime          `json:"updatedAt"`
}

// WebhookDetail :
type WebhookDetail struct {
	Webhook
	NoOfRetries int            `json:"noOfRetries"`
	Attempts    []WebhookRetry `json:"attempts"`
}

// WebhookRetry :
type WebhookRetry struct {
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	StatusCode  int               `json:"statusCode"`
	ElapsedTime int64             `json:"elapsedTime"`
	CreatedAt   DateTime          `json:"created"`
}
