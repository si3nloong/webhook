package shared

import (
	"log"
	"testing"

	"github.com/si3nloong/webhook/cmd"
)

func TestWebhook(t *testing.T) {
	svr := NewServer(&cmd.Config{})
	log.Println(svr)
}
