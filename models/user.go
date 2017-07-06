package models

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"log"

	gql "github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Password string     `json:"password"`
	Login    string     `json:"login" gorm:"unique_index"`
	Banned   bool       `json:"banned"`
	Comments []*Comment `json:"comments" gorm:"ForeignKey:AuthorID"`
}

type UserMutationResponse struct {
	Id uint `json:"id"`
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
			"comments": &gql.Field{
				Type:    gql.NewList(commentType),
				Resolve: resolveCommentsForUser,
			},
		},
	},
)

func userResolverSelect(p gql.ResolveParams) (interface{}, error) {
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
		return SelectUserByID(db, uint(id))
	} else {
		return nil, FieldNotFoundError("id")
	}
}

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

func resolverUpdate(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*Database)
	if db == nil {
		panic(errors.New("Can't find `Database` in context"))
	}
	l, ok := p.Args["login"]
	if !ok {
		return nil, FieldNotFoundError("login")
	}
	usr, err := SelectUserByLogin(db, l.(string))
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

	return UpdateUser(db, usr)
}

func resolverCreate(p gql.ResolveParams) (interface{}, error) {
	db := p.Context.Value("Database").(*Database)
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
	usr := User{
		Login:    l.(string),
		Password: pass.(string),
		Banned:   false,
	}
	return CreateUser(db, &usr)
}

func SelectUserByID(db *Database, id uint) (*User, error) {
	var u User
	if err := db.db.Where("id = ?", id).First(&u); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return &u, nil
}

func SelectUserByLogin(db *Database, login string) (*User, error) {
	var u User
	log.Println(login)
	if err := db.db.Where("login = ?", login).First(&u); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return &u, nil
}

func DBError(err error) error {
	return fmt.Errorf("Database error: %v", err)
}

func (u *User) EncryptPass() {
	h := sha1.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
}

func CreateUser(db *Database, u *User) (*UserMutationResponse, error) {
	u.EncryptPass()
	if err := db.db.Create(u); err.Error != nil {
		return &UserMutationResponse{0}, DBError(err.Error)
	} else {
		return &UserMutationResponse{u.ID}, nil
	}
}

func UpdateUser(db *Database, u *User) (*UserMutationResponse, error) {
	u.EncryptPass()
	if err := db.db.Save(u); err.Error != nil {
		return &UserMutationResponse{0}, DBError(err.Error)
	} else {
		return &UserMutationResponse{u.ID}, nil
	}
}
