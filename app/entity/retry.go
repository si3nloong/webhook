package entity

import "time"

type Retry struct {
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	StatusCode  int               `json:"statusCode"`
	ElapsedTime int64             `json:"elapsedTime"`
	CreatedAt   time.Time         `json:"createdAt"`
}
