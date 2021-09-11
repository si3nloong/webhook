package http

import (
	"fmt"
	"strings"

	"github.com/si3nloong/webhook/grpc/proto"
	"github.com/valyala/fasthttp"
)

func (s *Server) SendWebhook(ctx *fasthttp.RequestCtx) {
	var i struct {
		// CUrl    string            `json:"curl" validate:"omitempty,required"`
		URL     string            `json:"url" validate:"required,url,max=1000"`
		Method  string            `json:"method" validate:"oneof=GET POST"`
		Body    string            `json:"body" validate:"max=2048"`
		Headers map[string]string `json:"headers"`
	}

	i.URL = strings.TrimSpace(i.URL)
	i.Method = strings.TrimSpace(i.Method)

	if err := s.StructCtx(ctx, i); err != nil {
		return
	}

	req := new(proto.SendWebhookRequest)
	req.Url = i.URL
	switch i.Method {
	case "GET":
		req.Method = proto.SendWebhookRequest_GET
	case "POST":
		req.Method = proto.SendWebhookRequest_POST
	}
	req.Body = i.Body
	req.Headers = i.Headers

	// push to nats
	if err := s.Publish(ctx, req); err != nil {
		return
	}

	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}
