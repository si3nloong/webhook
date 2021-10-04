package dto

type Webhook struct {
	ID        string            `json:"id"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Body      string            `json:"body"`
	Headers   map[string]string `json:"headers"`
	Timeout   uint              `json:"timeout"`
	CreatedAt DateTime          `json:"createdAt"`
	UpdatedAt DateTime          `json:"updatedAt"`
}
