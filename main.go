package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fasthttp/router"
	validator "github.com/go-playground/validator/v10"
	"github.com/si3nloong/webhook/cmd"
	rpc "github.com/si3nloong/webhook/grpc"
	"github.com/si3nloong/webhook/grpc/proto"
	rest "github.com/si3nloong/webhook/http"
	"github.com/si3nloong/webhook/pubsub/nats"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(pwd)

	viper.SetEnvPrefix("WEBHOOK")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	log.Println("config =>", viper.ConfigFileUsed())

	// viper.SetConfigType("env")
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }

	v := validator.New()
	cfg := cmd.Config{}
	cfg.SetDefault()
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	// validate yaml value
	if err := v.StructCtx(ctx, cfg); err != nil {
		panic(err)
	}

	log.Println("cluster =>", cfg.MessageQueue.Redis.Cluster)

	// setup message queuing
	var (
		mq         = nats.New(&cfg)
		grpcServer *grpc.Server
	)

	if cfg.Enabled {
		go func() {
			svr := rest.NewServer(mq, v)
			httpServer := router.New()
			httpServer.GET("/health", svr.Health)
			httpServer.POST("/v1/webhook/send", svr.SendWebhook)
			log.Printf("RESTful serve at %s", cfg.Port)

			if err := fasthttp.ListenAndServe(":"+cfg.Port, httpServer.Handler); err != nil {
				defer cancel()
				panic(err)
			}
		}()
	}

	if cfg.GRPC.Enabled {
		grpcServer = grpc.NewServer()

		go func() {
			svr := rpc.NewServer(mq, v)
			proto.RegisterCurlHookServiceServer(grpcServer, svr)
			lis, err := net.Listen("tcp", ":"+cfg.GRPC.Port)
			if err != nil {
				panic(err)
			}

			log.Printf("gRPC serve at %s", cfg.GRPC.Port)
			if err := grpcServer.Serve(lis); err != nil {
				panic(err)
			}
		}()
	}

	select {
	case v := <-quit:
		log.Println("Quit", v)
		if grpcServer != nil {
			grpcServer.GracefulStop()
		}
		// s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Println("ctx.Done: %v", done)
	}

}
