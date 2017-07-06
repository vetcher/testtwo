package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"net/http"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	u "github.com/vetcher/testtwo/user"
	"log"
	"flag"
)

type DBConfig struct {
	Host *string
	User *string
	DBName *string
	Password *string
	Port *uint
}

const POSTGRES string = "postgres"

func parseDBConfig() *DBConfig {
	conf := DBConfig{
		User: flag.String("user", POSTGRES, "User"),
		Password: flag.String("password", POSTGRES, "User's password"),
		DBName: flag.String("db", POSTGRES, "Name of database"),
		Port: flag.Uint("port", 5432, "Postgres port"),
		Host: flag.String("host", POSTGRES, "Address of server"),
	}
	return &conf
}

func init() {
	conf := parseDBConfig()
	var err error
	u.MainDb, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			conf.Host,
			conf.Port,
			conf.User,
			conf.DBName,
			conf.Password,
		))
	if err != nil {
		panic(err)
	}
	u.MainDb.AutoMigrate(&u.User{})
}

func main() {
	mainschema, err := gql.NewSchema(
		gql.SchemaConfig{
			Query:    u.QueryType,
			Mutation: u.UserMutation,
		},
	)
	if err != nil {
		panic(err)
	}
	h := handler.New(&handler.Config{
		Schema: &mainschema,
		Pretty: true,
	})
	defer u.MainDb.Close()
	log.Println("Serve")
	http.ListenAndServe(":8080", h)
}
