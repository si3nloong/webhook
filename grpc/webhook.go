package grpc

import (
	"context"
	"log"
	"time"

	"github.com/si3nloong/curlhook/grpc/proto"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/status"
)

func (s *Server) SendWebhook(ctx context.Context, req *proto.SendWebhookRequest) (*proto.SendWebhookResponse, error) {
	if err := s.StructCtx(ctx, req); err != nil {
		return nil, status.Convert(err).Err()
	}

	// push to nats
	if err := s.mq.Publish(ctx); err != nil {
		return nil, status.Convert(err).Err()
	}

	go func() {
		log.Println("1")
		httpReq := fasthttp.AcquireRequest()
		httpResp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(httpReq)
		defer fasthttp.ReleaseResponse(httpResp)
		httpReq.Header.SetRequestURI(req.Url)
		httpReq.Header.SetMethod(req.Method.String())

		log.Println("2")

		for k, v := range req.Headers {
			httpReq.Header.Add(k, v)
		}
		httpReq.AppendBodyString(req.Body)

		// By default timeout is 5 seconds
		timeout := 5 * time.Second
		if req.Timeout > 0 {
			timeout = time.Second * time.Duration(req.Timeout)
		}
		if err := fasthttp.DoTimeout(httpReq, httpResp, timeout); err != nil {
			log.Println("Error =>", err)
			return
		}

		statusCode := httpResp.StatusCode()

		log.Println(string(httpResp.Body()))
		log.Println("StatusCode =>", statusCode)

		// 100 - 199
		if statusCode < fasthttp.StatusOK {
			// 200 - 399
		} else if statusCode >= fasthttp.StatusBadRequest {
		}
	}()

	resp := new(proto.SendWebhookResponse)
	return resp, nil
}
