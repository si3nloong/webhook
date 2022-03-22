package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/si3nloong/webhook/app/server/graphql/graph"
	"github.com/si3nloong/webhook/app/server/graphql/graph/generated"
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

	ws := shared.NewServer(&cfg)

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

	resolver := generated.Config{Resolvers: &graph.Resolver{WebhookServer: ws}}
	resolver.Directives.Validate = func(ctx context.Context, obj interface{}, next graphql.Resolver, rule string) (interface{}, error) {
		pc := graphql.GetPathContext(ctx)
		// graphql.GetFieldContext(ctx)
		if pc.Field == nil {
			return next(ctx)
		}

		v, ok := obj.(map[string]interface{})
		if !ok {
			return next(ctx)
		}

		if err := ws.VarCtx(ctx, v[*pc.Field], rule); err != nil {
			// log.Println(reflect.TypeOf(err))
			return nil, err
		}

		return next(ctx)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(resolver))
	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
