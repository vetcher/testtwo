package resolvers

import (
	"fmt"

	gql "github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	v "gopkg.in/go-playground/validator.v9"
)

var validator *v.Validate = v.New()

var RootMutation = gql.NewObject(gql.ObjectConfig{
	Name: "RootMutation",
	Fields: gql.Fields{
		"createUser": &gql.Field{
			Type: userType,
			Args: gql.FieldConfigArgument{
				"login": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"password": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
			},
			Resolve: resolveCreateUser,
		},
		"updateUser": &gql.Field{
			Type: userType,
			Args: gql.FieldConfigArgument{
				"login": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"password": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.String),
				},
				"banned": &gql.ArgumentConfig{
					Type: gql.Boolean,
				},
			},
			Resolve: resolveUpdateUser,
		},
		"postComment": &gql.Field{
			Type: commentType,
			Args: gql.FieldConfigArgument{
				"text": &gql.ArgumentConfig{
					Type:        gql.NewNonNull(gql.String),
					Description: "Message",
				},
				"login": &gql.ArgumentConfig{
					Type:        gql.NewNonNull(gql.String),
					Description: "Nickname (login) of author",
				},
			},
			Resolve: resolvePostComment,
		},
		"deleteComment": &gql.Field{
			Type: gql.Boolean,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type:        gql.NewNonNull(gql.ID),
					Description: "ID of message, which should be deleted",
				},
			},
			Resolve: resolveDeleteComment,
		},
	},
})

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
				Resolve: resolveGetUser,
			},
			"comment": &gql.Field{
				Type: commentType,
				Args: gql.FieldConfigArgument{
					"id": &gql.ArgumentConfig{
						Type: gql.ID,
					},
				},
				Resolve: resolveGetComment,
			},
		},
	},
)

// Decode from input to output and validate output
func DecodeAndValidate(input, output interface{}) error {
	err := mapstructure.WeakDecode(input, output)
	if err != nil {
		return fmt.Errorf("can't decode input: %v", err)
	}
	err = validator.Struct(output)
	if err != nil {
		return fmt.Errorf("validation falled: %v", err)
	}
	return nil
}
