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

	rpc "github.com/si3nloong/webhook/app/grpc"
	rest "github.com/si3nloong/webhook/app/http/api"
	"github.com/si3nloong/webhook/app/util"
	"github.com/si3nloong/webhook/cmd"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
)

func main() {

	var (
		// mq      pubsub.MessageQueue
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

	// switch cmd.MessageQueueEngine(cfg.MessageQueue.Engine) {
	// case cmd.MessageQueueEngineRedis:
	// 	mq, err = redis.New(cfg)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case cmd.MessageQueueEngineNats:
	// 	mq = nats.New(cfg)
	// default:
	// 	panic(fmt.Sprintf("unsupported message queue engine %q", cfg.MessageQueue.Engine))
	// }

	// serve HTTP
	if cfg.Enabled {
		go func() error {
			svr := rest.NewServer(v)
			httpServer := router.New()
			httpServer.GET("/health", svr.Health)
			httpServer.POST("/v1/webhook/send", svr.SendWebhook)
			log.Printf("HTTP/RESTful server serve at %v", cfg.Port)

			if err := fasthttp.ListenAndServe(util.FormatPort(cfg.Port), httpServer.Handler); err != nil {
				return err
			}

			return nil
		}()
	}

	// serve HTTP
	if cfg.Monitor.Enabled {
		go func() error {
			svr := rest.NewServer(v)
			httpServer := router.New()
			httpServer.GET("/health", svr.Health)
			httpServer.POST("/v1/webhook/send", svr.SendWebhook)
			log.Printf("Monitor server serve at %v", cfg.Monitor.Port)

			if err := fasthttp.ListenAndServe(util.FormatPort(cfg.Port), httpServer.Handler); err != nil {
				return err
			}

			return nil
		}()
	}

	// serve gRPC
	if cfg.GRPC.Enabled {
		grpcSvr = rpc.NewServer(cfg, v)

		go func() error {
			lis, err := net.Listen("tcp", util.FormatPort(cfg.GRPC.Port))
			if err != nil {
				return err
			}

			log.Printf("gRPC server serve at %v", cfg.GRPC.Port)
			if err := grpcSvr.Serve(lis); err != nil {
				return err
			}

			return nil
		}()
	}

	select {
	case <-quit:
		// close gRPC server if it's exists
		if grpcSvr != nil {
			grpcSvr.GracefulStop()
		}
		log.Println("Quit")

	case <-ctx.Done():
		log.Println("ctx.Done!")
	}

}
