package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/si3nloong/webhook/app/http/restful/dto"
	"github.com/si3nloong/webhook/app/http/restful/transformer"
	pb "github.com/si3nloong/webhook/protobuf"
	"github.com/valyala/fasthttp"
)

func (s *Server) listWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	datas, _, err := s.GetWebhooks(ctx, "", 100)
	if err != nil {
		return
	}

	items := new(dto.Items)
	for _, data := range datas {
		items.Items = append(items.Items, transformer.ToWebhook(data))
	}
	log.Println(r.URL.Host)
	items.Size = len(datas)
	items.Links.Self = r.URL.Host
	items.Links.Previous = ""
	items.Links.Self = ""
	writeJson(w, http.StatusOK, items)
}

func (s *Server) findWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := strings.TrimSpace(mux.Vars(r)["id"])
	if id == "" {
		return
	}

	data, err := s.FindWebhook(ctx, id)
	if err != nil {
		return
	}

	writeJson(w, http.StatusOK, &dto.Item{Item: transformer.ToWebhook(data)})
}

func (s *Server) sendWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var i struct {
		URL        string            `json:"url" validate:"required,url,max=1000"`
		Method     string            `json:"method" validate:"omitempty,oneof=GET POST PATCH PUT DELETE"`
		Body       string            `json:"body" validate:"max=2048"`
		Headers    map[string]string `json:"headers"`
		Retry      uint              `json:"retry" validate:"omitempty,required,max=10"`
		Concurrent uint              `json:"concurrent"`
		Timeout    int               `json:"timeout" validate:"omitempty,required,max=10000"`
	}

	// default value
	i.Method = "GET"
	i.Timeout = 1000 // 1 second

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

	req := new(pb.SendWebhookRequest)
	req.Url = i.URL

	switch i.Method {
	case fasthttp.MethodGet:
		req.Method = pb.SendWebhookRequest_GET
	case fasthttp.MethodPost:
		req.Method = pb.SendWebhookRequest_POST
	case fasthttp.MethodPatch:
		req.Method = pb.SendWebhookRequest_POST
	case fasthttp.MethodPut:
		req.Method = pb.SendWebhookRequest_POST
	case fasthttp.MethodDelete:
		req.Method = pb.SendWebhookRequest_POST
	}

	req.Body = i.Body
	req.Headers = i.Headers

	// push to nats
	if err := s.Publish(ctx, req); err != nil {
		return
	}

	writeJson(w, http.StatusOK, &dto.Item{Item: nil})
}
