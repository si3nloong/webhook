package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

// ErrorLog :
type ErrorLog struct {
	ID        ksuid.KSUID
	Method    string
	Error     string
	CreatedAt time.Time
}
