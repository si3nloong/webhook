package model

import (
	"time"

	"github.com/segmentio/ksuid"
)

// Log :
type Log struct {
	ID        ksuid.KSUID
	Method    string
	Error     string
	CreatedAt time.Time
}
