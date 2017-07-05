package main

import (
	"github.com/jinzhu/gorm"
	gql "github.com/graphql-go/graphql"
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

type User struct {
	gorm.Model
	Password string `json:"password"`
	Login string `json:"login" gorm:"unique_index"`
	Banned bool `json:"-"`
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

var QueryType = gql.NewObject(
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

var gqlreturnedId = gql.NewObject(gql.ObjectConfig{
	Name: "User id, which returns after creation",
	Fields: gql.Fields{
		"id": &gql.Field{
			Type: gql.ID,
		},
	},
})

var UserMutation = gql.NewObject(gql.ObjectConfig{
	Name: "User Mutation",
	Fields: gql.Fields{
		"createUser": &gql.Field{
			Type: gqlreturnedId,
			Args: gql.FieldConfigArgument{
				"login": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"password": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				obj, ok := p.Source.(*User)
				if !ok {
					// IDK which error should be here
					return nil, errors.New("I think its panic")
				}
				return CreateUser(obj)
			},
		},
		"updateUser": &gql.Field{
			Type: gqlreturnedId,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type: gql.ID,
				},
				"login": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"password": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				obj, ok := p.Source.(*User)
				if !ok {
					// IDK which error should be here
					return nil, errors.New("I think its panic")
				}
				return UpdateUser(obj)
			},
		},
	},
})

func SelectUser(id int) (*User, error) {
	var u User
	if err := db.Where("id = ?", id).First(&u); err.Error != nil {
		return nil, err.Error
	}
	return &u, nil
}

func (u *User) EncryptPass() {
	h := sha1.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
}

func CreateUser(u *User) (uint, error) {
	u.EncryptPass()
    if err := db.Create(u); err.Error != nil {
		return 0, err.Error
	} else {
		return u.ID, nil
	}
}

func UpdateUser(u *User) (uint, error) {
	u.EncryptPass()
	if err := db.Save(u); err.Error != nil {
		return 0, err.Error
	} else {
		return u.ID, nil
	}
}