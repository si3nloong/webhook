package entity

import "time"

type Retry struct {
	Body       string    `json:"body"`
	StatusCode uint      `json:"statusCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
