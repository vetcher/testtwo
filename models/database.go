package models

import (
	"flag"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

const POSTGRES string = "postgres"

type Database struct {
	db *gorm.DB
	config *DBConfig
}

func parseDBConfig() *DBConfig {
	conf := DBConfig{
		User:     flag.String("user", POSTGRES, "User"),
		Password: flag.String("password", POSTGRES, "User's password"),
		DBName:   flag.String("db", POSTGRES, "Name of database"),
		Port:     flag.Uint("port", 5432, "Postgres port"),
		Host:     flag.String("host", "localhost", "Address of server"),
	}
	flag.Parse()
	return &conf
}

type DBConfig struct {
	Host     *string
	User     *string
	DBName   *string
	Password *string
	Port     *uint
}

func InitDB() *Database {
	var db Database
	db.config = parseDBConfig()
	connectionParams := fmt.Sprintf("host=%v port=%d user=%v dbname=%v sslmode=disable password=%v",
		*db.config.Host,
		*db.config.Port,
		*db.config.User,
		*db.config.DBName,
		*db.config.Password,
	)
	var err error
	db.db, err = gorm.Open("postgres", connectionParams)
	if err != nil {
		panic(err)
	}
	log.Println("Connection:", connectionParams)
	db.db.AutoMigrate(&User{})
	db.db.AutoMigrate(&Comment{})
	return &db
}

func (db *Database) Close() {
	db.Close()
}
