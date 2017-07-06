package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vetcher/testtwo/models"
	r "github.com/vetcher/testtwo/resolvers"
)

func initHandler() *handler.Handler {
	schema, err := gql.NewSchema(
		gql.SchemaConfig{
			Query:    r.QueryType,
			Mutation: r.RootMutation,
		},
	)
	if err != nil {
		panic(err)
	}
	return handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})
}

func contextHandlerFunc(ctx context.Context, h *handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		h.ContextHandler(ctx, w, r)
		log.Println(r.URL.Path, time.Now().Sub(now))
	})
}

func main() {
	flag.Usage()
	db := models.InitDB()
	h := initHandler()
	ctx := context.WithValue(context.Background(), "Database", db)
	defer db.Close()
	log.Println("Serve")
	http.ListenAndServe(":8080", contextHandlerFunc(ctx, h))
}
