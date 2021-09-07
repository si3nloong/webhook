package main

import (
	"bytes"
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

	pwd, err := os.Getwd()
	viper.AddConfigPath(pwd)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	buf := new(bytes.Buffer)
	viper.ReadConfig(buf)
	log.Println(buf.String())
	log.Println(pwd, err)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v := validator.New()
	pbc := nats.New()

	log.Println("redis cluster =>", viper.GetBool("CURLHOOK_REDIS_CLUSTER"))
	log.Println("redis host =>", viper.Get("CURLHOOK_REDIS_HOST"))
	log.Println("redis port =>", viper.Get("CURLHOOK_REDIS_PORT"))
	log.Println("redis password =>", viper.Get("CURLHOOK_REDIS_PASSWORD"))

	go func() {
		httpPort := "8000"
		svr := rest.NewServer(v)
		httpServer := router.New()
		httpServer.GET("/health", svr.Health)
		log.Printf("http serve at %s", httpPort)
		if err := fasthttp.ListenAndServe(":"+httpPort, httpServer.Handler); err != nil {
			log.Println(err)
			cancel()
		}
	}()

	cmd.Execute()

	grpcServer := grpc.NewServer()
	proto.RegisterCurlHookServiceServer(grpcServer, rpc.NewServer(v, pbc))

	go func() {
		grpcPort := "9000"
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			panic(err)
		}

		log.Printf("gRPC serve at %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	select {
	case v := <-quit:
		log.Println("Quit", v)
		grpcServer.GracefulStop()
		// s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Println("ctx.Done: %v", done)
	}
	log.Println(viper.Get("author"))
	// viper.SetConfigType("yaml")

}
