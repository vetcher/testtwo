package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"net/http"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var mainschema, _ = gql.NewSchema(
	gql.SchemaConfig{
		Query:    QueryType,
		Mutation: UserMutation,
	},
)

func resolverFunc(p gql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if ok {
		return SelectUser(id)
	}
	return nil, nil
}

var db gorm.DB

func init() {
	host := "localhost"
	user := "postgres"
	dbname := "omg_test"
	password := "postgres"
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			host,
			user,
			dbname,
			password,
		))
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
}

func main() {
	h := handler.New(&handler.Config{
		Schema: &mainschema,
		Pretty: true,
	})
	defer db.Close()
	http.ListenAndServe(":8080", h)
}
