package shared

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/avast/retry-go/v3"
	pb "github.com/si3nloong/webhook/grpc/proto"

	"github.com/valyala/fasthttp"
)

type WebhookServer struct {
}

func (ws *WebhookServer) SendWebhook(req *pb.SendWebhookRequest) error {
	err := retry.Do(
		func() error {
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
				return retry.Unrecoverable(err)
			} else if err != nil {
				// return retry.IsRecoverable(err)
				return err
			}

			log.Println("Response =======>")
			log.Println(httpResp.String())
			statusCode := httpResp.StatusCode()

			// 100 - 199
			if statusCode < fasthttp.StatusOK {
				// 200 - 399
			} else if statusCode >= fasthttp.StatusBadRequest {
				return errors.New("")
			}

			return nil
		},
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			log.Println("backoff retrying")
			log.Println(n, err, config)
			// 			fmt.Println("Server fails with: " + err.Error())
			// 			// if retriable, ok := err.(*retry.RetriableError); ok {
			// 			// 	fmt.Printf("Client follows server recommendation to retry after %v\n", retriable.RetryAfter)
			// 			// 	return retriable.RetryAfter
			// 			// }
			// 			// apply a default exponential back off strategy
			return retry.BackOffDelay(n, err, config)
		}),
		retry.Attempts(uint(req.Retry)),
	)
	if err != nil {
		return err
	}

	return nil
}
