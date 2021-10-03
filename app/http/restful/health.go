package http

import (
	"fmt"
	"net/http"
)

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, server is ok")
}
