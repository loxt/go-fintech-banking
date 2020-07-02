package database

import (
	"github.com/jinzhu/gorm"
	"github.com/loxt/go-fintech-banking/helpers"
)

var DB *gorm.DB

func InitDatabase() {
	db, err := gorm.Open(
		"postgres",
		"host=127.0.0.1 port=5432 user=postgres dbname=bankapp password=postgres sslmode=disable")

	helpers.HandleErr(err)

	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(200)

	DB = db
}
