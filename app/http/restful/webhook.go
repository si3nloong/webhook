package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/si3nloong/webhook/app/grpc/proto"
	"github.com/valyala/fasthttp"
)

func (s *Server) SendWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var i struct {
		URL        string            `json:"url" validate:"required,url,max=1000"`
		Method     string            `json:"method" validate:"oneof=GET POST PATCH PUT DELETE"`
		Body       string            `json:"body" validate:"max=2048"`
		Headers    map[string]string `json:"headers"`
		Retry      uint              `json:"retry" validate:"omitempty,required,max=10"`
		Concurrent uint              `json:"concurrent"`
		Timeout    int               `json:"timeout" validate:"omitempty,required,max=10000"`
	}

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		log.Println("Error =>", err)
		return
	}

	i.URL = strings.TrimSpace(i.URL)
	i.Method = strings.TrimSpace(i.Method)
	i.Body = strings.TrimSpace(i.Body)

	if err := s.Validate(i); err != nil {
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
		req.Method = proto.SendWebhookRequest_POST
	case fasthttp.MethodPut:
		req.Method = proto.SendWebhookRequest_POST
	case fasthttp.MethodDelete:
		req.Method = proto.SendWebhookRequest_POST
	}

	req.Body = i.Body
	req.Headers = i.Headers

	// push to nats
	if err := s.Publish(ctx, req); err != nil {
		return
	}

	fmt.Fprintf(w, "Hello, %s!\n", "xxx")
}
