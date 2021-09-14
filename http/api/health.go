package http

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func (s *Server) Health(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}
