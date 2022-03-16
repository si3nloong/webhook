package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/si3nloong/webhook/app/shared"
)

type Server struct {
	*mux.Router
	shared.WebhookServer
}

func NewServer(ws shared.WebhookServer) *Server {
	svr := new(Server)
	svr.WebhookServer = ws
	svr.Router = mux.NewRouter()

	svr.Router.HandleFunc("/health", svr.health)
	svr.Router.HandleFunc("/v1/webhooks", svr.listWebhooks).Methods("GET")
	svr.Router.HandleFunc("/v1/webhook/{id}", svr.findWebhook).Methods("GET")
	svr.Router.HandleFunc("/v1/webhook/send", svr.sendWebhook).Methods("POST")
	return svr
}

func writeJson(w http.ResponseWriter, statusCode int, dest interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(dest)
	if err != nil {
		log.Println(err)
	}
	w.Write(b)
}
