package shared

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/si3nloong/webhook/cmd"
	es "github.com/si3nloong/webhook/database/elasticsearch"
	"github.com/si3nloong/webhook/entity"
	pb "github.com/si3nloong/webhook/grpc/proto"
	"github.com/si3nloong/webhook/mq/nats"
	"github.com/si3nloong/webhook/mq/redis"
	"github.com/valyala/fasthttp"
)

/*
Flow --------->
send webhook ---> insert into Queue
get from queue ---> fire webhook
fire webhook ---> record stat (add success count or log error)
*/

type Repository interface {
	InsertLog(ctx context.Context, data *entity.Log) error
	GetLogs(ctx context.Context) ([]entity.Log, error)
	FindLog(ctx context.Context, id string) (*entity.Log, error)

	// GetStats(ctx context.Context) error
	// IncrStat(ctx context.Context) error
}

type WebhookServer interface {
	Validate(ctx context.Context, src interface{}) error
	SendWebhook(ctx context.Context, req *pb.SendWebhookRequest) error
	Repository

	// pubsub.MessageQueue
	// db.Repository
}

type webhookServer struct {
	v *validator.Validate
	Repository
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
		nats.New(cfg)
	case cmd.MessageQueueEngineRedis:
		redis.New(cfg)
	default:
		panic(fmt.Sprintf("invalid database engine %s", cfg.DB.Engine))
	}
	if err != nil {
		panic(err)
	}

	log.Println(svr.Repository)
	// svr.Stat = mt.NewMetricServerWithRedisClient(redis.NewClient(&redis.Options{}))
	// svr.db, err = sql.Open("mysql", "root:abcd1234@/webhook")
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = svr.db.Exec("CREATE TABLE `ErrorLog` (`ID` VARCHAR(50), `Method` VARCHAR(50), `Error` TEXT);")
	// log.Println(err)
	return svr
}

func (s *webhookServer) Validate(ctx context.Context, src interface{}) error {
	return s.v.StructCtx(ctx, src)
}

func (s *webhookServer) SendWebhook(ctx context.Context, req *pb.SendWebhookRequest) error {
	opts := make([]retry.Option, 0)
	if req.Retry < 1 {
		req.Retry = 1
	}

	opts = append(opts, retry.Attempts(uint(req.Retry)))
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
				return err
			}

			log.Println("Response =======>")
			log.Println(httpResp.String())
			statusCode := httpResp.StatusCode()

			// 100 - 199
			if statusCode < fasthttp.StatusOK {
				log.Println("100 - 199")
				// 500
			} else if statusCode >= fasthttp.StatusInternalServerError {
				log.Println("500")
				return &requestError{body: httpResp.String()}
				// 400
			} else if statusCode >= fasthttp.StatusBadRequest {
				log.Println("400")
			}

			return nil
		},
		opts...,
	)
	if err != nil {
		log.Println("Error here =>", err)
		// s.LogError(ctx,, req, err)
		return err
	}

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
