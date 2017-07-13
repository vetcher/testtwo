package resolvers

import (
	"fmt"

	gql "github.com/graphql-go/graphql"
	"github.com/vetcher/testtwo/models"
)

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
			"comments": &gql.Field{
				Type:    gql.NewList(commentType),
				Resolve: resolveCommentsForUser,
			},
		},
	},
)

func resolveGetUser(p gql.ResolveParams) (interface{}, error) {
	var id OnlyID
	err := DecodeAndValidate(p.Args, &id)
	if err != nil {
		return nil, err
	}
	db := p.Context.Value("Database").(*models.Database)
	return models.SelectUserByID(db, id.ID)
}

func resolveUpdateUser(p gql.ResolveParams) (interface{}, error) {
	var u User
	err := DecodeAndValidate(p.Args, &u)
	if err != nil {
		return nil, err
	}
	db := p.Context.Value("Database").(*models.Database)
	usr, err := models.SelectUserByLogin(db, u.Login)
	if err != nil {
		return nil, fmt.Errorf("can't find `User` by login: %v", err)
	}
	usr.Password = u.Password
	usr.IsBanned = u.IsBanned

	return models.UpdateUser(db, usr)
}

func resolveCreateUser(p gql.ResolveParams) (interface{}, error) {
	var u User
	err := DecodeAndValidate(p.Args, &u)
	if err != nil {
		return nil, err
	}
	usr := models.User{
		Login:    u.Login,
		Password: u.Password,
		IsBanned: u.IsBanned,
	}
	db := p.Context.Value("Database").(*models.Database)
	return models.CreateUser(db, &usr)
}
