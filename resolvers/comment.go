package resolvers

import (
	"errors"
	"fmt"

	gql "github.com/graphql-go/graphql"
	"github.com/vetcher/testtwo/models"
)

var commentType = gql.NewObject(gql.ObjectConfig{
	Name:        "Comment",
	Description: "User's comment",
	Fields: gql.Fields{
		"text": &gql.Field{
			Type:        gql.String,
			Description: "User's message",
		},
	},
})

func init() {
	commentType.AddFieldConfig("author", &gql.Field{
		Type:        userType,
		Description: "Author of this comment",
	})
}

func resolvePostComment(p gql.ResolveParams) (interface{}, error) {
	var in Comment
	if err := DecodeAndValidate(p.Args, &in); err != nil {
		return nil, err
	}
	db := p.Context.Value("Database").(*models.Database)
	u, err := models.SelectUserByLogin(db, in.Login)
	if err != nil {
		return nil, fmt.Errorf("can't find user becuse of: %v", err)
	}
	c := models.Comment{
		AuthorID: u.ID,
		Text:     in.Text,
	}
	return models.PostComment(db, &c)
}

func resolveDeleteComment(p gql.ResolveParams) (interface{}, error) {
	var id OnlyID
	if err := DecodeAndValidate(p.Args, &id); err != nil {
		return false, err
	}
	db := p.Context.Value("Database").(*models.Database)
	return models.DeleteComment(db, id.ID)
}

func commentResolverSelect(p gql.ResolveParams) (interface{}, error) {
	var id OnlyID
	if err := DecodeAndValidate(p.Args, &id); err != nil {
		return nil, err
	}
	db := p.Context.Value("Database").(*models.Database)
	return models.SelectCommentByID(db, id.ID)
}

func resolveCommentsForUser(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*models.Database)
	if u, ok := p.Source.(*models.User); !ok {
		return nil, errors.New("input source not a `*User`")
	} else {
		return models.LoadCommentsForUser(db, u)
	}
}
