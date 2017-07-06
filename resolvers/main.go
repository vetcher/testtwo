package resolvers

import (
	"fmt"

	gql "github.com/graphql-go/graphql"
)

var RootMutation = gql.NewObject(gql.ObjectConfig{
	Name: "RootMutation",
	Fields: gql.Fields{
		"createUser": &gql.Field{
			Type: gqlReturnedID,
			Args: gql.FieldConfigArgument{
				"login": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"password": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: resolverCreate,
		},
		"updateUser": &gql.Field{
			Type: gqlReturnedID,
			Args: gql.FieldConfigArgument{
				"login": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"password": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"banned": &gql.ArgumentConfig{
					Type: gql.Boolean,
				},
			},
			Resolve: resolverUpdate,
		},
		"postComment": &gql.Field{
			Type: gqlReturnedID,
			Args: gql.FieldConfigArgument{
				"text": &gql.ArgumentConfig{
					Type:        gql.String,
					Description: "Message",
				},
				"login": &gql.ArgumentConfig{
					Type:        gql.String,
					Description: "Nickname (login) of author",
				},
			},
			Resolve: resolvePostComment,
		},
		"deleteComment": &gql.Field{
			Type: gqlReturnedID,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type:        gql.ID,
					Description: "ID of message, which should be deleted",
				},
			},
			Resolve: resolveDeleteComment,
		},
	},
})

func FieldNotFoundError(field string) error {
	return fmt.Errorf("Field `%v` not found empty", field)
}
