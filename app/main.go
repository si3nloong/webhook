package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	rpc "github.com/si3nloong/webhook/app/grpc"
	restful "github.com/si3nloong/webhook/app/http/restful"
	"github.com/si3nloong/webhook/app/shared"
	"github.com/si3nloong/webhook/app/util"
	"github.com/si3nloong/webhook/cmd"
	"github.com/spf13/viper"
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

	ws := shared.NewServer(cfg)

	// serve HTTP
	if cfg.Enabled {
		go func() error {
			r := mux.NewRouter()
			svr := restful.NewServer(ws)
			r.HandleFunc("/", svr.Health)
			r.HandleFunc("/v1/webhook/send", svr.SendWebhook).Methods("POST")

			log.Printf("HTTP server serve at %v", cfg.Port)
			log.Fatal(http.ListenAndServe(util.FormatPort(cfg.Port), r))
			return nil
		}()
	}

	// serve gRPC
	if cfg.GRPC.Enabled {
		grpcSvr = rpc.NewServer(cfg, ws)

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
