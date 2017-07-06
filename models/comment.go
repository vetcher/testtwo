package models

import (
	"errors"
	"fmt"
	"strconv"

	"log"

	gql "github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Text     string `json:"text"`
	AuthorID uint   `gorm:"index" json:"author"`
}

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

type CommentMutationResponse struct {
	Id uint `json:"id"`
}

func resolvePostComment(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	log.Println("OK")
	l := p.Args["login"]
	u, err := SelectUserByLogin(db, l.(string))
	if err != nil {
		return nil, fmt.Errorf("can't find user becuse of: %v", err)
	}
	text := p.Args["text"]
	c := Comment{
		AuthorID: u.ID,
		Text:     text.(string),
	}
	return PostComment(db, &c)
}

func PostComment(db *Database, c *Comment) (*CommentMutationResponse, error) {
	if err := db.db.Create(c); err.Error != nil {
		return &CommentMutationResponse{0}, DBError(err.Error)
	} else {
		return &CommentMutationResponse{c.ID}, nil
	}
}

func resolveDeleteComment(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*Database)
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
	return DeleteComment(db, uint(id))
}

func DeleteComment(db *Database, id uint) (*CommentMutationResponse, error) {
	var temp Comment
	if err := db.db.Where("id = ?", id).First(&temp); err.Error != nil {
		return &CommentMutationResponse{0}, DBError(err.Error)
	} else {
		if err := db.db.Delete(temp); err.Error != nil {
			return &CommentMutationResponse{0}, DBError(err.Error)
		}
		return &CommentMutationResponse{id}, nil
	}
}

func commentResolverSelect(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	idstr, ok := p.Args["id"]
	if ok {
		id, err := strconv.ParseUint(idstr.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Parsing `id` error: %v", err)
		}
		return SelectCommentByID(db, uint(id))
	} else {
		return nil, FieldNotFoundError("id")
	}
}

func SelectCommentByID(db *Database, id uint) (*Comment, error) {
	var c Comment
	if err := db.db.Where("id = ?", id).First(&c); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return &c, nil
}

func LoadCommentsForUser(db *Database, u *User) ([]*Comment, error) {
	if err := db.db.Model(u).Related(&u.Comments, "AuthorID"); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return u.Comments, nil
}

func resolveCommentsForUser(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	if u, ok := p.Source.(*User); !ok {
		return nil, errors.New("Not a `*User`")
	} else {
		return LoadCommentsForUser(db, u)
	}
}
