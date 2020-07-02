package main

import (
	"github.com/loxt/go-fintech-banking/api"
	"github.com/loxt/go-fintech-banking/database"
)

func main() {
	database.InitDatabase()
	api.StartApi()
}
