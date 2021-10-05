package shared

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/avast/retry-go/v3"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	es "github.com/si3nloong/webhook/app/database/elasticsearch"
	"github.com/si3nloong/webhook/app/entity"
	"github.com/si3nloong/webhook/app/mq/nats"
	"github.com/si3nloong/webhook/app/mq/redis"
	"github.com/si3nloong/webhook/cmd"
	pb "github.com/si3nloong/webhook/protobuf"
	"github.com/valyala/fasthttp"
)

/*
Flow --------->
send webhook ---> insert into Queue
get from queue ---> fire webhook
fire webhook ---> record stat (add success count or log error)
*/

type Repository interface {
	CreateWebhook(ctx context.Context, data *entity.WebhookRequest) error
	GetWebhooks(ctx context.Context, curCursor string, limit uint) (datas []*entity.WebhookRequest, nextCursor string, err error)
	FindWebhook(ctx context.Context, id string) (*entity.WebhookRequest, error)
	UpdateWebhook(ctx context.Context, id string, status int, body string) error
}

type MessageQueue interface {
	Publish(ctx context.Context, data *entity.WebhookRequest) error
}

type WebhookServer interface {
	Repository
	Validate(src interface{}) error
	Publish(ctx context.Context, req *pb.SendWebhookRequest) (*entity.WebhookRequest, error)
}

type webhookServer struct {
	// logger log.Logger
	Repository
	v  *validator.Validate
	mq MessageQueue
}

func NewServer(cfg cmd.Config) WebhookServer {
	var (
		svr = new(webhookServer)
		err error
	)

	svr.v = validator.New()

	// setup Database
	switch cfg.DB.Engine {
	case cmd.DatabaseEngineElasticsearch:
		svr.Repository, err = es.New(cfg)
	default:
		panic(fmt.Sprintf("invalid database engine %s", cfg.DB.Engine))
	}
	if err != nil {
		panic(err)
	}

	// setup Message Queueing
	switch cfg.MessageQueue.Engine {
	case cmd.MessageQueueEngineNSQ:
	case cmd.MessageQueueEngineNats:
		svr.mq, err = nats.New(cfg)
	case cmd.MessageQueueEngineRedis:
		svr.mq, err = redis.New(cfg, func(delivery rmq.Delivery) {
			data := entity.WebhookRequest{}
			if err := json.Unmarshal([]byte(delivery.Payload()), &data); err != nil {
				return
			}

			svr.fireWebhook(&data)
		})
	default:
		panic(fmt.Sprintf("invalid database engine %s", cfg.DB.Engine))
	}
	if err != nil {
		panic(err)
	}

	return svr
}

func (s *webhookServer) Validate(src interface{}) error {
	return s.v.Struct(src)
}

func (s *webhookServer) Publish(ctx context.Context, req *pb.SendWebhookRequest) (*entity.WebhookRequest, error) {
	utcNow := time.Now().UTC()

	// may be we store into database first before publish to message queue
	data := entity.WebhookRequest{}
	data.ID = ksuid.New()
	data.Method = req.Method.String()
	data.URL = req.Url
	data.Headers = make(map[string]string)
	for k, v := range req.Headers {
		data.Headers[k] = v
	}
	data.Body = req.Body
	data.Timeout = 3000 // 3 seconds
	if req.Timeout > 0 {
		data.Timeout = uint(req.Timeout)
	}
	data.Retries = make([]entity.Retry, 0)
	data.CreatedAt = utcNow
	data.UpdatedAt = utcNow

	if err := s.CreateWebhook(ctx, &data); err != nil {
		return nil, err
	}

	if err := s.mq.Publish(ctx, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *webhookServer) fireWebhook(data *entity.WebhookRequest) error {
	ctx := context.TODO()
	// opts := make([]retry.Option, 0)
	// if req.Retry < 1 {
	// 	req.Retry = 1
	// }
	log.Println("SendWebhook now....!!!")

	// opts = append(opts, retry.Attempts(uint(req.Retry)))
	// if req.RetryMechanism > 0 {
	// 	retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
	// 		log.Println("backoff retrying")
	// 		log.Println(n, err, config)
	// 		// fmt.Println("Server fails with: " + err.Error())
	// 		// if retriable, ok := err.(*retry.RetriableError); ok {
	// 		// 	fmt.Printf("Client follows server recommendation to retry after %v\n", retriable.RetryAfter)
	// 		// 	return retriable.RetryAfter
	// 		// }

	// 		// apply a default exponential back off strategy
	// 		return retry.BackOffDelay(n, err, config)
	// 	})
	// }

	// err := retry.Do(
	// 	func() error {
	httpReq := fasthttp.AcquireRequest()
	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(httpReq)
	defer fasthttp.ReleaseResponse(httpResp)
	httpReq.Header.SetRequestURI(data.URL)
	httpReq.Header.SetMethod(data.Method)

	for k, v := range data.Headers {
		httpReq.Header.Add(k, v)
	}
	httpReq.AppendBodyString(data.Body)

	// By default timeout is 3 seconds
	timeout := 3 * time.Second
	// if data.Timeout > 0 {
	// 	timeout = time.Second * time.Duration(req.Timeout)
	// }

	log.Println("Request =======>")
	log.Println(httpReq.String())

	var dnsError *net.DNSError
	if err := fasthttp.DoTimeout(httpReq, httpResp, timeout); errors.As(err, &dnsError) {
		// If it's an invalid host, drop the request directly
		log.Println("Error 1 =======>", err)
		return retry.Unrecoverable(err)
	} else if err != nil {
		switch t := err.(type) {
		case *net.OpError:
			// if it's an unknown host, drop it
			if t.Op == "dial" {
				println("Unknown host")
			} else if t.Op == "read" {
				println("Connection refused")
			}

		case syscall.Errno:
			if t == syscall.ECONNREFUSED {
				println("Connection refused")
			}
		}
		log.Println("Error 2 =======>", err, reflect.TypeOf(err))
		return err
	}

	log.Println("Response =======>")
	log.Println(httpResp.String())
	statusCode := httpResp.StatusCode()
	i64, _ := strconv.ParseInt(string(httpResp.Header.Peek("Content-Length")), 10, 64)
	if i64 > 1048 {
		// discard the response if body bigger than 1mb
	}
	body := httpResp.Body()

	// 100 - 199
	if statusCode < fasthttp.StatusOK {
		log.Println("100 - 199")
		// 500
	} else if statusCode >= fasthttp.StatusInternalServerError {
		log.Println("500")
		// return &requestError{body: httpResp.String()}
		// 400
	} else if statusCode >= fasthttp.StatusBadRequest {
		log.Println("400")
	}

	s.UpdateWebhook(ctx, data.ID.String(), statusCode, string(body))

	// s.Update(ctx, id)
	// 		return nil
	// 	},
	// 	opts...,
	// )
	// if err != nil {
	// 	log.Println("Error here =>", err)
	// 	// s.LogError(ctx,, req, err)
	// 	return err
	// }

	// if err := s.Incr(ctx, metric.StatTypeSucceed); err != nil {
	// 	return err
	// }
	return nil
}

type requestError struct {
	body string
}

func (e requestError) Error() string {
	return ""
}
