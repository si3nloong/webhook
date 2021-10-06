package http

import (
	"encoding/json"
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
	items.Items = make([]interface{}, 0)
	for _, data := range datas {
		items.Items = append(items.Items, transformer.ToWebhook(data))
	}
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

	writeJson(w, http.StatusOK, &dto.Item{Item: transformer.ToWebhookDetail(data)})
}

func (s *Server) sendWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var i struct {
		URL     string            `json:"url" validate:"required,url,max=1000"`
		Method  string            `json:"method" validate:"omitempty,oneof=GET POST PATCH PUT DELETE"`
		Body    string            `json:"body" validate:"max=2048"`
		Headers map[string]string `json:"headers"`
		Retry   struct {
			Max      uint8  `json:"max" validate:"omitempty,required,max=10"`
			Strategy string `json:"strategy" validate:"omitempty,required,oneof=backoff"`
		} `json:"retry"`
		Concurrent uint8 `json:"concurrent"`
		Timeout    uint  `json:"timeout" validate:"omitempty,required,max=10000"`
	}

	// default value
	i.Method = "GET"
	i.Retry.Max = 10
	i.Retry.Strategy = "backoff"
	i.Timeout = 1000 // 1 second

	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		writeJson(w, http.StatusBadRequest, err)
		return
	}

	i.URL = strings.TrimSpace(i.URL)
	i.Method = strings.TrimSpace(i.Method)
	i.Body = strings.TrimSpace(i.Body)

	if err := s.Validate(i); err != nil {
		writeJson(w, http.StatusUnprocessableEntity, err)
		return
	}

	req := pb.SendWebhookRequest{}

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

	req.Url = i.URL
	req.Body = i.Body
	req.Headers = i.Headers
	req.Retry = uint32(i.Retry.Max)

	// push to nats
	data, err := s.Publish(ctx, &req)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, http.StatusOK, &dto.Item{Item: transformer.ToWebhook(data)})
}
