package resolvers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

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
	db := p.Context.Value("Database").(*models.Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	log.Println("OK")
	l := p.Args["login"]
	u, err := models.SelectUserByLogin(db, l.(string))
	if err != nil {
		return nil, fmt.Errorf("can't find user becuse of: %v", err)
	}
	text := p.Args["text"]
	c := models.Comment{
		AuthorID: u.ID,
		Text:     text.(string),
	}
	return models.PostComment(db, &c)
}

func resolveDeleteComment(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*models.Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	idstr, ok := p.Args["id"].(string)
	if !ok {
		return nil, FieldNotFoundError("id")
	}
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse `%v` to uint", idstr)
	}
	return models.DeleteComment(db, uint(id))
}

func commentResolverSelect(p gql.ResolveParams) (interface{}, error) {
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
		return models.SelectCommentByID(db, uint(id))
	} else {
		return nil, FieldNotFoundError("id")
	}
}

func resolveCommentsForUser(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*models.Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	if u, ok := p.Source.(*models.User); !ok {
		return nil, errors.New("Not a `*User`")
	} else {
		return models.LoadCommentsForUser(db, u)
	}
}
