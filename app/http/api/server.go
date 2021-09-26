package http

import (
	"github.com/si3nloong/webhook/app/shared"
)

type Server struct {
	ws shared.WebhookServer
}

func NewServer(ws shared.WebhookServer) *Server {
	svr := new(Server)
	return svr
}
