package http

import (
	"fmt"
	"strings"

	"github.com/si3nloong/webhook/app/grpc/proto"
	"github.com/valyala/fasthttp"
)

func (s *Server) SendWebhook(ctx *fasthttp.RequestCtx) {
	var i struct {
		URL     string            `json:"url" validate:"required,url,max=1000"`
		Method  string            `json:"method" validate:"oneof=GET POST PATCH PUT DELETE"`
		Body    string            `json:"body" validate:"max=2048"`
		Headers map[string]string `json:"headers"`
	}

	i.URL = strings.TrimSpace(i.URL)
	i.Method = strings.TrimSpace(i.Method)

	if err := s.Validate(ctx, i); err != nil {
		return
	}

	req := new(proto.SendWebhookRequest)
	req.Url = i.URL
	switch i.Method {
	case fasthttp.MethodGet:
		req.Method = proto.SendWebhookRequest_GET
	case fasthttp.MethodPost:
		req.Method = proto.SendWebhookRequest_POST
	case fasthttp.MethodPatch:
	case fasthttp.MethodPut:
	case fasthttp.MethodDelete:
	}
	req.Body = i.Body
	req.Headers = i.Headers

	// push to nats
	if err := s.Publish(ctx, req); err != nil {
		return
	}

	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}
