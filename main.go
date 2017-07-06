package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	m "github.com/vetcher/testtwo/models"
)

const POSTGRES string = "postgres"

type DBConfig struct {
	Host     *string
	User     *string
	DBName   *string
	Password *string
	Port     *uint
}

func parseDBConfig() *DBConfig {
	conf := DBConfig{
		User:     flag.String("user", POSTGRES, "User"),
		Password: flag.String("password", POSTGRES, "User's password"),
		DBName:   flag.String("db", POSTGRES, "Name of database"),
		Port:     flag.Uint("port", 5432, "Postgres port"),
		Host:     flag.String("host", "localhost", "Address of server"),
	}
	flag.Parse()
	return &conf
}

func initDB() *gorm.DB {
	conf := parseDBConfig()
	connectionParams := fmt.Sprintf("host=%v port=%d user=%v dbname=%v sslmode=disable password=%v",
		*conf.Host,
		*conf.Port,
		*conf.User,
		*conf.DBName,
		*conf.Password,
	)
	db, err := gorm.Open("postgres", connectionParams)
	if err != nil {
		panic(err)
	}
	log.Println("Connection:", connectionParams)
	db.AutoMigrate(&m.User{})
	db.AutoMigrate(&m.Comment{})
	return db
}

func initHandler() *handler.Handler {
	schema, err := gql.NewSchema(
		gql.SchemaConfig{
			Query:    m.QueryType,
			Mutation: m.UserMutation,
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
		h.ContextHandler(ctx, w, r)
	})
}

func main() {
	flag.Usage()
	db := initDB()
	h := initHandler()
	ctx := context.WithValue(context.Background(), "Database", db)
	defer db.Close()
	log.Println("Serve")
	http.ListenAndServe(":8080", contextHandlerFunc(ctx, h))
}
