package resolvers

import (
	"errors"
	"fmt"
	"strconv"

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
				Resolve: userResolverSelect,
			},
			"comment": &gql.Field{
				Type: commentType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{
						Type: gql.ID,
					},
				},
				Resolve: commentResolverSelect,
			},
		},
	},
)

var gqlReturnedID = gql.NewObject(gql.ObjectConfig{
	Name: "InstanceID",
	Fields: gql.Fields{
		"id": &gql.Field{
			Type: gql.ID,
		},
	},
})

func userResolverSelect(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*models.Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	idstr, ok := p.Args["id"]
	if ok {
		id, err := strconv.ParseUint(idstr.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Parsing `id` error: %v", err)
		}
		return models.SelectUserByID(db, uint(id))
	} else {
		return nil, FieldNotFoundError("id")
	}
}

func resolverUpdate(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*models.Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	l, ok := p.Args["login"]
	if !ok {
		return nil, FieldNotFoundError("login")
	}
	usr, err := models.SelectUserByLogin(db, l.(string))
	if err != nil {
		return nil, err
	}
	pass, ok := p.Args["password"]
	if !ok {
		return nil, FieldNotFoundError("password")
	}
	ban, ok := p.Args["banned"]
	if !ok {
		return nil, FieldNotFoundError("banned")
	}
	usr.Password = pass.(string)
	usr.Banned = ban.(bool)

	return models.UpdateUser(db, usr)
}

func resolverCreate(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*models.Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	l, ok := p.Args["login"]
	if !ok {
		return nil, FieldNotFoundError("login")
	}
	pass, ok := p.Args["password"]
	if !ok {
		return nil, FieldNotFoundError("login")
	}
	usr := models.User{
		Login:    l.(string),
		Password: pass.(string),
		Banned:   false,
	}
	return models.CreateUser(db, &usr)
}
