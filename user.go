package main

import (
	"time"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Password string
	Login string
	RegistrationDate time.Time
	Banned bool
}
