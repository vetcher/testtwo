package main

import (
	"github.com/jinzhu/gorm"
	"github.com/graphql-go/graphql"
)

type User struct {
	gorm.Model
	Password string `json:"password"`
	Login string `json:"login"`
	Banned bool `json:"banned"`
}

func (u *User) GraphqlFields()graphql.Fields {
	return graphql.Fields{
		"login": &graphql.Field{
			Type: graphql.String,
		},
		"password": &graphql.Field{
			Type: graphql.String,
		},
		"banned": &graphql.Field{
			Type: graphql.Boolean,
		},
	}
}