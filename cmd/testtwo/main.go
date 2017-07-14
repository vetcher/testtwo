package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"flag"

	"fmt"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	client "github.com/vetcher/comments-msv/transport"
	"github.com/vetcher/testtwo/models"
	r "github.com/vetcher/testtwo/resolvers"
	"google.golang.org/grpc"
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
		begin := time.Now()
		h.ContextHandler(ctx, w, r)
		defer log.Println(r.URL.Path, time.Since(begin))
	})
}

type CommentSvcConfig struct {
	Host *string
	Port *uint
}

func NewCommentSvcConfig() *CommentSvcConfig {
	return &CommentSvcConfig{
		Port: flag.Uint("commentport", 10000, "CommentSVC port"),
		Host: flag.String("commenthost", "localhost", "Address of CommentSVC server"),
	}
}

func main() {
	DBConf := models.NewDBConfig()
	commentSvcConf := NewCommentSvcConfig()
	flag.Parse()
	db := models.InitDB(DBConf)
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *commentSvcConf.Host, *commentSvcConf.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	commentSvc := client.NewClient(conn)
	h := initHandler()
	ctx := context.WithValue(context.Background(), "Database", db)
	ctx = context.WithValue(ctx, "CommentService", commentSvc)
	defer db.Close()
	log.Println("Serve :8080")
	http.ListenAndServe(":8080", contextHandlerFunc(ctx, h))
}
