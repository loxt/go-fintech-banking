package main

import (
	"github.com/loxt/go-fintech-banking/api"
	"github.com/loxt/go-fintech-banking/migrations"
)

func main() {
	migrations.MigrateTransactions()
	//migrations.Migrate()
	api.StartApi()
}
