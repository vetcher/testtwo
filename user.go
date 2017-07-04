package main

import (
	"github.com/jinzhu/gorm"
	gql "github.com/graphql-go/graphql"
)

type User struct {
	gorm.Model
	Password string `json:"password"`
	Login string `json:"login"`
	Banned bool `json:"banned"`
}

var userType = gql.NewObject(
	gql.ObjectConfig{
		Name: "User",
		Fields: gql.Fields{
			"login": &gql.Field{
				Type: gql.String,
			},
			"password": &gql.Field{
				Type: gql.String,
			},
			"banned": &gql.Field{
				Type: gql.Boolean,
			},
		},
	},
)

var queryType = gql.NewObject(
	gql.ObjectConfig{
		Name: "Query",
		Fields: gql.Fields{
			"user": &gql.Field{
				Type: userType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{
						Type: gql.ID,
					},
				},
				Resolve: resolverFunc,
			},
		},
	},
)

func SelectUser(id int) (*User, error) {
	var u User
	if err := db.Where("id = ?", id).First(&u); err != nil {
		return nil, err.Error
	}
	return &u, nil
}

func CreateUser(u *User) error {
	
	return nil
}