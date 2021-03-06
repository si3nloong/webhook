package shared

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/adjust/rmq/v4"
	"github.com/avast/retry-go/v3"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
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

type webhookServer struct {
	Repository
	v  *validator.Validate
	mq MessageQueue
}

func NewServer(cfg *cmd.Config) WebhookServer {
	var (
		err error
		svr = &webhookServer{
			v: validator.New(),
		}
	)

	// setup Database
	switch cfg.DB.Engine {
	case cmd.DatabaseEngineElasticsearch:
	// 	svr.Repository, err = es.New(cfg)
	default:
		// 	panic(fmt.Sprintf("invalid database engine %s", cfg.DB.Engine))
	}
	if err != nil {
		panic(err)
	}

	log.Println("Database engine =>", cfg.DB.Engine)
	log.Println("Message queue engine =>", cfg.MessageQueue.Engine)

	// setup Message Queueing
	switch cfg.MessageQueue.Engine {
	case cmd.MessageQueueEngineRedis:
		{
			svr.mq, err = redis.New(cfg, func(msg rmq.Delivery) {
				var (
					data = entity.WebhookRequest{}
					errs error
				)

				// capture error if exists when it's end
				defer func() {
					if errs != nil {
						svr.logErrorIfAny(errs)
						svr.logErrorIfAny(msg.Reject())
						return
					}

					svr.logErrorIfAny(msg.Ack())
				}()

				if errs = json.Unmarshal([]byte(msg.Payload()), &data); errs != nil {
					return
				}

				if errs = svr.fireWebhook(&data); errs != nil {
					return
				}
			})
		}
	case cmd.MessageQueueEngineNSQ:
		{
		}
	case cmd.MessageQueueEngineNats:
		{
			svr.mq, err = nats.New(cfg)
		}
	default:
		// by default it will select in-memory pubsub
	}
	if err != nil {
		panic(err)
	}

	return svr
}

func (s *webhookServer) Validate(src interface{}) error {
	return s.v.Struct(src)
}

func (s *webhookServer) VarCtx(ctx context.Context, src interface{}, tag string) error {
	return s.v.VarCtx(ctx, src, tag)
}

func (s *webhookServer) logErrorIfAny(err error) {
	if err != nil {
		s.LogError(err)
	}
}

func (*webhookServer) LogError(err error) {
	log.Println("Error", err)
}

func (s *webhookServer) Publish(ctx context.Context, req *pb.SendWebhookRequest) (*entity.WebhookRequest, error) {
	utcNow := time.Now().UTC()

	// Store the request to DB first before publishing it to message queue
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
	data.Attempts = make([]entity.Attempt, 0)
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
	startTime := time.Now().UTC()
	opts := make([]retry.Option, 0)

	log.Println("SendWebhook now....!!!")

	opts = append(opts, retry.Attempts(10))
	// if req.RetryMechanism > 0 {
	opts = append(opts, retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
		log.Println("backoff retrying")
		log.Println(n, err, config)
		// fmt.Println("Server fails with: " + err.Error())
		// if retriable, ok := err.(*retry.RetriableError); ok {
		// 	fmt.Printf("Client follows server recommendation to retry after %v\n", retriable.RetryAfter)
		// 	return retriable.RetryAfter
		// }

		// apply a default exponential back off strategy
		return retry.BackOffDelay(n, err, config)
	}))
	// }

	errs := retry.Do(
		func() error {
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
			var body string
			i64, _ := strconv.ParseInt(string(httpResp.Header.Peek("Content-Length")), 10, 64)
			// discard the response if body bigger than 1mb
			if i64 < 1048 {
				body = string(httpResp.Body())
			}

			utcNow := time.Now().UTC()
			att := entity.Attempt{}
			att.Headers = make(map[string]string)
			httpResp.Header.VisitAll(func(key, value []byte) {
				att.Headers[b2s(key)] = b2s(value)
			})
			att.Body = body
			att.ElapsedTime = time.Now().UTC().Sub(startTime).Milliseconds()
			att.StatusCode = uint(statusCode)
			att.CreatedAt = utcNow

			if err := s.UpdateWebhook(ctx, data.ID.String(), &att); err != nil {
				return err
			}

			// 100 - 199
			if statusCode < fasthttp.StatusOK {
				log.Println("100 - 199")
				// 500
			} else if statusCode >= fasthttp.StatusInternalServerError {
				log.Println("500")
				// return &requestError{body: httpResp.String()}
				// 400
				return errors.New("error 500")
			} else if statusCode >= fasthttp.StatusBadRequest {
				log.Println("400")
				// retry
				return errors.New("error 404")
			}

			return nil
		},
		opts...,
	)

	return errs
}

type requestError struct {
	body string
}

func (e requestError) Error() string {
	return ""
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
