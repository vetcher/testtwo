package models

import (
	"errors"
	"fmt"
	"strconv"

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

func commentResolverSelect(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database")
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	idstr, ok := p.Args["id"]
	if ok {
		id, err := strconv.ParseUint(idstr.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Parsing `id` error: %v", err)
		}
		return SelectCommentByID(db.(*gorm.DB), uint(id))
	} else {
		return nil, FieldNotFoundError("id")
	}
}

func SelectCommentByID(db *gorm.DB, id uint) (*Comment, error) {
	var c Comment
	if err := db.Where("id = ?", id).First(&c); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return &c, nil
}

func LoadCommentsForUser(db *gorm.DB, u *User) ([]*Comment, error) {
	if err := db.Model(u).Related(&u.Comments, "AuthorID"); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return u.Comments, nil
}

func resolveCommentsForUser(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database")
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	if u, ok := p.Source.(*User); !ok {
		return nil, errors.New("Not a `*User`")
	} else {
		return LoadCommentsForUser(db.(*gorm.DB), u)
	}
}
