package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/si3nloong/webhook/app/shared"
)

type Server struct {
	shared.Server
}

func NewServer(v *validator.Validate) *Server {
	svr := new(Server)
	// svr.Validate = v
	// svr.MessageQueue = mq
	return svr
}
