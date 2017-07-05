package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"net/http"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	u "testtwo/user"
	"log"
)

var mainschema, ttt = gql.NewSchema(
	gql.SchemaConfig{
		Query:    u.QueryType,
		Mutation: u.UserMutation,
	},
)

func init() {
	host := "localhost"
	user := "postgres"
	dbname := "omg_test"
	password := "postgres"
	var err error
	u.MainDb, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			host,
			user,
			dbname,
			password,
		))
	if err != nil {
		panic(err)
	}
	u.MainDb.AutoMigrate(&u.User{})
}

func main() {
	log.Println(ttt)
	h := handler.New(&handler.Config{
		Schema: &mainschema,
		Pretty: true,
	})
	defer u.MainDb.Close()
	log.Println("Serve")
	http.ListenAndServe(":8080", h)
}
