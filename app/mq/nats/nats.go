package nats

import (
	"errors"
	"log"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/nats-io/nats.go"
	pb "github.com/si3nloong/webhook/app/grpc/proto"
	"github.com/si3nloong/webhook/cmd"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
)

type natsMQ struct {
	sync.RWMutex
	subj string
	js   nats.JetStreamContext
	subs []*nats.Subscription
}

func New(cfg cmd.Config) (*natsMQ, error) {
	q := new(natsMQ)

	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	// if err := handleMessage(&pb.SendWebhookRequest{
	// 	Url: "https://183jkashdkjhasjkdh.com",
	// }); err != nil {
	// 	panic(err)
	// }

	log.Println(nc.Statistics)

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		panic(err)
	}

	streamName := "webhook"
	stream, err := js.StreamInfo(streamName)
	log.Println(stream, err)
	// js.DeleteStream(streamName)
	if err == nats.ErrStreamNotFound {
		if _, err := js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"test"},
		}); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	for i := 0; i < cfg.NoOfWorker; i++ {
		q.subs = append(q.subs)
	}

	{
		sub, err := js.QueueSubscribe(
			"test",
			"webhook1",
			q.onQueueSubscribe,
			nats.ManualAck(),
			nats.AckWait(10*time.Second),
			nats.MaxDeliver(10),
		)

		log.Println("Subscription =>", sub, err)
	}

	q.js = js

	return q, nil
}

func (mq *natsMQ) onQueueSubscribe(msg *nats.Msg) {
	log.Println("Handle message ========>")
	log.Println(string(msg.Data))
	// log.Println(msg)
	req := new(pb.SendWebhookRequest)

	if err := proto.Unmarshal(msg.Data, req); err != nil {
		log.Println(err)
		return
	}

	if err := handleMessage(req); err != nil {
		log.Println(err)
		// msg.Nak()
		return
	}

	// if everything ok, acknowledge and don't retry
	msg.Ack()
}

func handleMessage(req *pb.SendWebhookRequest) error {
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

func (q *natsMQ) GracefulStop() error {
	// return q.
	return nil
}
