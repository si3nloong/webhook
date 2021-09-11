package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fasthttp/router"
	validator "github.com/go-playground/validator/v10"

	"github.com/si3nloong/webhook/cmd"
	rpc "github.com/si3nloong/webhook/grpc"
	rest "github.com/si3nloong/webhook/http"
	"github.com/si3nloong/webhook/pubsub"
	"github.com/si3nloong/webhook/pubsub/nats"
	"github.com/si3nloong/webhook/pubsub/redis"
	"github.com/si3nloong/webhook/util"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
)

func main() {

	var (
		mq      pubsub.MessageQueue
		grpcSvr *grpc.Server
		quit    = make(chan os.Signal, 1)
		v       = validator.New()
		cfg     = cmd.Config{}
	)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get current path
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.AddConfigPath(pwd)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("webhook")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// log.Println("config =>", viper.ConfigFileUsed())

	// set the default value for configuration
	cfg.SetDefault()
	// read config into struct
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	// validate yaml value
	if err := v.StructCtx(ctx, cfg); err != nil {
		panic(err)
	}

	// setup message queuing
	switch cmd.MessageQueueEngine(cfg.MessageQueue.Engine) {
	case cmd.MessageQueueEngineRedis:
		mq = redis.New(cfg)
	case cmd.MessageQueueEngineNats:
		mq = nats.New(cfg)
	// case cmd.MessageQueueEngineNSQ:
	// 	mq = redis.New(ctx, cfg)
	default:
	}

	// serve http
	if cfg.Enabled {
		go func() {
			svr := rest.NewServer(mq, v)
			httpServer := router.New()
			httpServer.GET("/health", svr.Health)
			httpServer.POST("/v1/webhook/send", svr.SendWebhook)
			log.Printf("HTTP/RESTful server serve at %v", cfg.Port)

			if err := fasthttp.ListenAndServe(util.FormatPort(cfg.Port), httpServer.Handler); err != nil {
				defer cancel()
				panic(err)
			}
		}()
	}

	// serve gRPC
	if cfg.GRPC.Enabled {
		grpcSvr = rpc.NewServer(cfg, mq, v)

		go func() {
			lis, err := net.Listen("tcp", util.FormatPort(cfg.GRPC.Port))
			if err != nil {
				panic(err)
			}

			log.Printf("gRPC server serve at %v", cfg.GRPC.Port)
			if err := grpcSvr.Serve(lis); err != nil {
				panic(err)
			}
		}()
	}

	select {
	case <-quit:
		// close gRPC server if it's exists
		if grpcSvr != nil {
			grpcSvr.GracefulStop()
		}

	case done := <-ctx.Done():
		log.Println("ctx.Done: ", done)
	}

}
