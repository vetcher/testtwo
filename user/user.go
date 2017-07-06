package user

import (
	"github.com/jinzhu/gorm"
	gql "github.com/graphql-go/graphql"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strconv"
	"fmt"
)

var MainDb *gorm.DB

type User struct {
	gorm.Model
	Password string `json:"password"`
	Login string `json:"login" gorm:"unique_index"`
	Banned bool `json:"-"`
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
		},
	},
)


func resolverFunc(p gql.ResolveParams) (interface{}, error) {
	idstr, ok := p.Args["id"]
	if ok {
		id, err := strconv.ParseUint(idstr.(string), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Parsing `id` error: %v", err)
		}
		return SelectUser(uint(id))
	}
	return nil, fmt.Errorf("Field `id` not found")
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
				Resolve: resolverFunc,
			},
		},
	},
)

var gqlReturnedID = gql.NewObject(gql.ObjectConfig{
	Name: "UserId",
	Fields: gql.Fields{
		"id": &gql.Field{
			Type: gql.ID,
		},
	},
})

var UserMutation = gql.NewObject(gql.ObjectConfig{
	Name: "UserMutation",
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
	},
})

func FieldNotFoundError(field string) error {
	return fmt.Errorf("Field `%v` not found empty", field)
}

func resolverUpdate(p gql.ResolveParams) (interface{}, error) {
	l, ok := p.Args["login"]
	if !ok {
		return nil, FieldNotFoundError("login")
	}
	usr, err := SelectUserByLogin(l.(string))
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
	return UpdateUser(usr)
}

func resolverCreate(p gql.ResolveParams) (interface{}, error) {
	l, ok := p.Args["login"]
	if !ok {
		return nil, errors.New("I think its panic")
	}
	pass, ok := p.Args["password"]
	if !ok {
		return nil, errors.New("I think its panic")
	}
	usr := User{
		Login: l.(string),
		Password: pass.(string),
		Banned: false,
	}
	return CreateUser(&usr)
}

func SelectUser(id uint) (*User, error) {
	var u User
	if err := MainDb.Where("id = ?", id).First(&u); err.Error != nil {
		return nil, fmt.Errorf("Database error: %v", err.Error)
	}
	return &u, nil
}

func SelectUserByLogin(login string) (*User, error) {
	var u User
	if err := MainDb.Where("login = ?", login).First(&u); err.Error != nil {
		return nil, fmt.Errorf("Database error: %v", err.Error)
	}
	return &u, nil
}

func (u *User) EncryptPass() {
	h := sha1.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
}

func CreateUser(u *User) (*UserMutationResponse, error) {
	u.EncryptPass()
    if err := MainDb.Create(u); err.Error != nil {
		return &UserMutationResponse{0}, fmt.Errorf("Database error: %v", err.Error)
	} else {
		return &UserMutationResponse{u.ID}, nil
	}
}

func UpdateUser(u *User) (*UserMutationResponse, error) {
	u.EncryptPass()
	if err := MainDb.Save(u); err.Error != nil {
		return &UserMutationResponse{0}, fmt.Errorf("Database error: %v", err.Error)
	} else {
		return &UserMutationResponse{u.ID}, nil
	}
}