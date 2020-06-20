package useraccounts

import (
	"github.com/loxt/go-fintech-banking/helpers"
	"github.com/loxt/go-fintech-banking/interfaces"
)

func updateAccount(id uint, amount int) {
	db := helpers.ConnectDB()
	db.Model(&interfaces.Account{}).
		Where("id = ?", id).
		Update("balance", amount)
	defer db.Close()
}
