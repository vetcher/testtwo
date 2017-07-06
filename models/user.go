package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"

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

func DBError(err error) error {
	return fmt.Errorf("Database error: %v", err)
}

func (u *User) EncryptPass() {
	h := sha1.New()
	h.Write([]byte(u.Password))
	u.Password = hex.EncodeToString(h.Sum(nil))
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
