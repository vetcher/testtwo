package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
)

func init() {
	host := "localhost"
	user := "postgres"
	dbname := "omg_test"
	password := "postgres"
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
			host,
			user,
			dbname,
			password,
		))
	if err != nil {
		panic(err)
	}
	defer db.Close()

}

func main() {
}
