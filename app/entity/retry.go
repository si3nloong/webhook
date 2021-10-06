package entity

import "time"

type Retry struct {
	Response struct {
		Body       string `json:"body"`
		StatusCode int    `json:"statusCode"`
	} `json:"response"`
	CreatedAt time.Time `json:"createdAt"`
}
