package http

import (
	"encoding/json"
	"net/http"
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJson(w, http.StatusOK, json.RawMessage(`{"message":"ok"}`))
}
