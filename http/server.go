package http

import (
	"github.com/go-playground/validator/v10"
)

type Server struct {
	*validator.Validate
}

func NewServer(v *validator.Validate) *Server {
	return &Server{
		Validate: v,
	}
}
