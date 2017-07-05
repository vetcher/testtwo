package user

import (
	"github.com/jinzhu/gorm"
	gql "github.com/graphql-go/graphql"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"log"
	"strconv"
)

var MainDb *gorm.DB

type User struct {
	gorm.Model
	Password string `json:"password"`
	Login string `json:"login" gorm:"unique_index"`
	Banned bool `json:"-"`
}

type Answer struct {
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
	log.Println("ID is", idstr)
	if ok {
		id, err := strconv.ParseUint(idstr.(string), 10, 64)
		if err != nil {
			return nil, err
		}
		return SelectUser(uint(id))
	}
	return nil, nil
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

var gqlreturnedId = gql.NewObject(gql.ObjectConfig{
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
			Type: gqlreturnedId,
			Args: gql.FieldConfigArgument{
				"login": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"password": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				l, ok := p.Args["login"]
				if !ok {
					return nil, errors.New("I think its panic 1")
				}
				pass, ok := p.Args["password"]
				if !ok {
					return nil, errors.New("I think its panic 2")
				}
				usr := User{
					Login: l.(string),
					Password: pass.(string),
					Banned: false,
				}
				return CreateUser(&usr)
			},
		},
		"updateUser": &gql.Field{
			Type: gqlreturnedId,
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

func resolverUpdate(p gql.ResolveParams) (interface{}, error) {
	l, ok := p.Args["login"]
	if !ok {
		return nil, errors.New("Login is empty")
	}
	usr, err := SelectUserByLogin(l.(string))
	if err != nil {
		return nil, err
	}
	pass, ok := p.Args["password"]
	if !ok {
		return nil, errors.New("I think its panic")
	}
	ban, ok := p.Args["banned"]
	//b, err := strconv.ParseBool(ban.(string))
	if !ok || err != nil {
		return nil, errors.New("I think its panic")
	}
	usr.Password = pass.(string)
	usr.Banned = ban.(bool)
	return UpdateUser(usr)
}

func SelectUser(id uint) (*User, error) {
	var u User
	if err := MainDb.Where("id = ?", id).First(&u); err.Error != nil {
		return nil, err.Error
	}
	return &u, nil
}

func SelectUserByLogin(login string) (*User, error) {
	var u User
	if err := MainDb.Where("login = ?", login).First(&u); err.Error != nil {
		return nil, err.Error
	}
	return &u, nil
}

func (u *User) EncryptPass() {
	h := sha1.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
}

func CreateUser(u *User) (Answer, error) {
	u.EncryptPass()
    if err := MainDb.Create(u); err.Error != nil {
		return Answer{0}, err.Error
	} else {
		return Answer{u.ID}, nil
	}
}

func UpdateUser(u *User) (Answer, error) {
	u.EncryptPass()
	if err := MainDb.Save(u); err.Error != nil {
		return Answer{0}, err.Error
	} else {
		return Answer{u.ID}, nil
	}
}