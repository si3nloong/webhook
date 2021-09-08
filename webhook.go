package main

import (
	"errors"
	"log"
	"net"
	"reflect"
	"time"

	"github.com/avast/retry-go"
	pb "github.com/si3nloong/webhook/grpc/proto"
	"github.com/valyala/fasthttp"
)

func sendWebhook(req *pb.SendWebhookRequest) error {
	retry.Do(
		func() error {
			return nil
		},
		retry.Attempts(3),
	)
	httpReq := fasthttp.AcquireRequest()
	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(httpReq)
	defer fasthttp.ReleaseResponse(httpResp)
	httpReq.Header.SetRequestURI(req.Url)
	httpReq.Header.SetMethod(req.Method.String())

	for k, v := range req.Headers {
		httpReq.Header.Add(k, v)
	}
	httpReq.AppendBodyString(req.Body)

	// By default timeout is 5 seconds
	timeout := 5 * time.Second
	if req.Timeout > 0 {
		timeout = time.Second * time.Duration(req.Timeout)
	}

	log.Println("Request =======>")
	log.Println(httpReq.String())

	var dnsError *net.DNSError
	if err := fasthttp.DoTimeout(httpReq, httpResp, timeout); errors.As(err, &dnsError) {
		// If it's a invalid host, drop the request directly
		log.Println(err, reflect.TypeOf(err))
		return err
	} else if err != nil {
		return err
	}

	log.Println("Response =======>")
	log.Println(httpResp.String())
	statusCode := httpResp.StatusCode()

	// 100 - 199
	if statusCode < fasthttp.StatusOK {
		// 200 - 399
	} else if statusCode >= fasthttp.StatusBadRequest {
	}

	return nil
}
