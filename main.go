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

	viper.AddConfigPath(pwd)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config := cmd.Config{}
	config.Enabled = true
	config.Port = "3000"
	config.GRPC.Port = "9000"
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	v := validator.New()
	pbc := nats.New()

	if config.Enabled {
		go func() {
			svr := rest.NewServer(v)
			httpServer := router.New()
			httpServer.GET("/health", svr.Health)
			log.Printf("RESTful serve at %s", config.Port)

			if err := fasthttp.ListenAndServe(":"+config.Port, httpServer.Handler); err != nil {
				defer cancel()
				panic(err)
			}
		}()
	}

	grpcServer := grpc.NewServer()

	if config.GRPC.Enabled {
		proto.RegisterCurlHookServiceServer(grpcServer, rpc.NewServer(v, pbc))
		go func() {
			lis, err := net.Listen("tcp", ":"+config.GRPC.Port)
			if err != nil {
				panic(err)
			}

			log.Printf("gRPC serve at %s", config.GRPC.Port)
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalf("failed to serve: %s", err)
			}
		}()
	}

	select {
	case v := <-quit:
		log.Println("Quit", v)
		grpcServer.GracefulStop()
		// s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Println("ctx.Done: %v", done)
	}

}
