package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
	"encoding/json"
	"net/http"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"errors"
	"io/ioutil"
	"log"
)

var getschema, _ = gql.NewSchema(
	gql.SchemaConfig{
		Query: queryType,
		Mutation: userMutation,
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

func handlerfunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		response := gql.Do(gql.Params{
			Schema:        getschema,
			RequestString: r.URL.Query()["query"][0],
		})
		json.NewEncoder(w).Encode(response)
	} else if r.Method == "POST" {
		// TODO
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		} else {
			response := gql.Do(gql.Params{
				Schema:        postschema,
				RequestString: string(body),
			})
			json.NewEncoder(w).Encode(response)
		}
	} else {
		json.NewEncoder(w).Encode(
			gqlerrors.FormatError(
				errors.New("Wrong method, only GET and POST supported"),
		))
	}
}

func main() {
	http.HandleFunc("/graphql", handlerfunc)
	defer db.Close()
	http.ListenAndServe(":8080", nil)
}
