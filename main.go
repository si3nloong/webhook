package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/si3nloong/curlhook/cmd"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, %s!\n", ctx.UserValue("name"))
}

func main() {

	cmd.Execute()
	log.Println(viper.Get("author"))
	// viper.SetConfigType("yaml")
	r := router.New()
	r.GET("/", Hello)
	// r.GET("/hello/{name}", Hello)
	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))

	// lis, err := net.Listen("tcp", ":9000")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// grpcServer := grpc.NewServer()

	// if err := grpcServer.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %s", err)
	// }
}
