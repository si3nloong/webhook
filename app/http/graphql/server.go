package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/si3nloong/webhook/app/http/graphql/graph"
	"github.com/si3nloong/webhook/app/http/graphql/graph/generated"
	"github.com/si3nloong/webhook/app/shared"
	"github.com/si3nloong/webhook/cmd"
	"github.com/spf13/viper"
)

const defaultPort = "8080"

func main() {
	var cfg cmd.Config
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

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ws := shared.NewServer(cfg)

	r := mux.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{
			"Apollographql-Client-Name",
			"Apollographql-Client-Version",
			"Content-Type",
		},
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{WebhookServer: ws}}))

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
