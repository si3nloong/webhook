package http

import (
	"github.com/si3nloong/webhook/app/shared"
)

type Server struct {
	shared.WebhookServer
}

func NewServer(ws shared.WebhookServer) *Server {
	svr := new(Server)
	svr.WebhookServer = ws
	return svr
}
