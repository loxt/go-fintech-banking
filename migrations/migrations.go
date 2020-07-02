package migrations

import (
	"github.com/loxt/go-fintech-banking/database"
	"github.com/loxt/go-fintech-banking/helpers"
	"github.com/loxt/go-fintech-banking/interfaces"
)

func createAccounts() {
	users := &[2]interfaces.User{
		{
			Username: "Loxt",
			Email:    "loxt@loxt.com",
		},
		{
			Username: "gio",
			Email:    "gio@gio.com",
		},
	}

	for i := 0; i < len(users); i++ {
		generatePassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := &interfaces.User{
			Username: users[i].Username,
			Email:    users[i].Email,
			Password: generatePassword,
		}
		database.DB.Create(&user)

		account := &interfaces.Account{
			Type:    "Daily Account",
			Name:    users[i].Username + "'s account",
			Balance: uint(10000*i + 1),
			UserID:  user.ID,
		}
		database.DB.Create(&account)
	}
	defer database.DB.Close()
}

func Migrate() {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	Transactions := &interfaces.Transaction{}
	database.DB.AutoMigrate(&User, &Account, &Transactions)

	createAccounts()
}
