package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"

	rpc "github.com/si3nloong/webhook/app/server/grpc"
	rest "github.com/si3nloong/webhook/app/server/rest"
	"github.com/si3nloong/webhook/app/shared"
	"github.com/si3nloong/webhook/app/util"
	"github.com/si3nloong/webhook/cmd"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {

	var (
		grpcSvr *grpc.Server
		quit    = make(chan os.Signal, 1)
		v       = validator.New()
		wsCfg   = new(cmd.Config)
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

	// set the default value for configuration
	wsCfg.SetDefault()
	// read config into struct
	if err := viper.Unmarshal(&wsCfg); err != nil {
		panic(err)
	}

	// validate yaml value
	if err := v.StructCtx(ctx, wsCfg); err != nil {
		panic(err)
	}

	ws := shared.NewServer(wsCfg)

	// serve HTTP
	if wsCfg.Enabled {
		go func() error {
			log.Printf("HTTP server serve at %v", wsCfg.Port)
			log.Fatal(http.ListenAndServe(util.FormatPort(wsCfg.Port), rest.NewServer(ws)))
			return nil
		}()
	}

	// serve gRPC
	if wsCfg.GRPC.Enabled {
		grpcSvr = rpc.NewServer(wsCfg, ws)

		go func() error {
			lis, err := net.Listen("tcp", util.FormatPort(wsCfg.GRPC.Port))
			if err != nil {
				return err
			}

			log.Printf("gRPC server serve at %v", wsCfg.GRPC.Port)
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
