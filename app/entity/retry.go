package entity

import "time"

type Retry struct {
	StatusCode uint      `json:"statusCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
