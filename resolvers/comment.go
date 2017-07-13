package resolvers

import (
	"errors"
	"fmt"

	"time"

	gql "github.com/graphql-go/graphql"
	commentsvc "github.com/vetcher/comments-msv/models"
	"github.com/vetcher/comments-msv/service"
	"github.com/vetcher/testtwo/models"
)

var commentType = gql.NewObject(gql.ObjectConfig{
	Name:        "Comment",
	Description: "User's comment",
	Fields: gql.Fields{
		"Text": &gql.Field{
			Type:        gql.String,
			Description: "User's message",
		},
		"created_at": &gql.Field{
			Type:        gql.String,
			Name:        "CreatedAt",
			Description: "Server-side time when comment was created (posted)",
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				return p.Source.(*commentsvc.Comment).CreatedAt.Format(time.Stamp), nil
			},
		},
		// more fields adds in runtime, please check `init()`
	},
})

func init() {
	commentType.AddFieldConfig("author", &gql.Field{
		Type:        userType,
		Description: "Author of this comment",
		Resolve: func(p gql.ResolveParams) (interface{}, error) {
			comment := p.Source.(*commentsvc.Comment)
			db := p.Context.Value("Database").(*models.Database)
			return models.SelectUserByID(db, comment.AuthorID)
		},
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
	client := p.Context.Value("CommentService").(service.CommentService)
	return client.PostComment(u.ID, in.Text)
}

func resolveDeleteComment(p gql.ResolveParams) (interface{}, error) {
	var id OnlyID
	if err := DecodeAndValidate(p.Args, &id); err != nil {
		return false, err
	}
	client := p.Context.Value("CommentService").(service.CommentService)
	return client.DeleteCommentByID(id.ID)
}

func resolveGetComment(p gql.ResolveParams) (interface{}, error) {
	var id OnlyID
	if err := DecodeAndValidate(p.Args, &id); err != nil {
		return nil, err
	}
	client := p.Context.Value("CommentService").(service.CommentService)
	comment, err := client.GetCommentByID(id.ID)
	return comment, err
}

func resolveCommentsForUser(p gql.ResolveParams) (interface{}, error) {
	if u, ok := p.Source.(*models.User); !ok {
		return nil, errors.New("input source not a `*User`")
	} else {
		client := p.Context.Value("CommentService").(service.CommentService)
		return client.GetCommentsByAuthorID(u.ID)
	}
}
