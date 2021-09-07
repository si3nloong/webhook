package nats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	pb "github.com/si3nloong/webhook/grpc/proto"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	subj string
	js   nats.JetStreamContext
}

func New() *Client {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

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

	{
		sub, err := js.QueueSubscribe(
			"test",
			"webhook1",
			func(msg *nats.Msg) {
				log.Println("Handle message ========>")
				log.Println(string(msg.Data))
				log.Println(msg)
				req := new(pb.SendWebhookRequest)

				if err := proto.Unmarshal(msg.Data, req); err != nil {
					log.Println(err)
					return
				}

				if err := handleMessage(req); err != nil {
					log.Println(err)
					msg.Nak()
					return
				}

				msg.Ack()
			},
			nats.AckExplicit(),
			// nats.Durable("webhook"),
			// nats.AckWait(30*time.Second),
			nats.MaxDeliver(10),
		)

		log.Println("Subscription =>", sub, err)
	}

	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		ack, err := js.Publish("test", []byte("hello world!"))
	// 		log.Println(ack, err)
	// 	}
	// }()

	// stan.AckWait(20 * time.Second)
	// sc.QueueSubscribe(
	// 	"webhook",
	// 	"webhook",
	// 	func(msg *stan.Msg) {},
	// 	stan.AckWait(5*time.Second),
	// 	stan.DurableName("webhook"),
	// 	stan.MaxInflight(5),
	// )
	return &Client{
		js:   js,
		subj: "test",
	}
}

func handleMessage(req *pb.SendWebhookRequest) error {
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

	if err := fasthttp.DoTimeout(httpReq, httpResp, timeout); err != nil {
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
