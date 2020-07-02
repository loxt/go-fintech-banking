package transactions

import (
	"github.com/loxt/go-fintech-banking/database"
	"github.com/loxt/go-fintech-banking/helpers"
	"github.com/loxt/go-fintech-banking/interfaces"
)

func CreateTransaction(From uint, To uint, Amount int) {
	transaction := &interfaces.Transaction{From: From, To: To, Amount: Amount}
	database.DB.Create(&transaction)
}

//func MigrateTransactions() {
//	Transactions := &interfaces.Transaction{}
//
//	database.DB.AutoMigrate(&Transactions)
//	defer database.DB.Close()
//}

func GetTransactionsByAccount(id uint) []interfaces.ResponseTransaction {
	var transactions []interfaces.ResponseTransaction
	database.DB.Table("transactions").
		Select("id, transactions.from, transactions.to, amount").
		Where(interfaces.Transaction{From: id}).Or(interfaces.Transaction{To: id}).
		Scan(&transactions)

	return transactions
}

func GetMyTransactions(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)

	if isValid {
		var accounts []interfaces.ResponseAccount
		database.DB.Table("accounts").Select("id, name, balance").
			Where("user_id = ?", id).Scan(&accounts)

		var transactions []interfaces.ResponseTransaction

		for i := 0; i < len(accounts); i++ {
			accTransactions := GetTransactionsByAccount(accounts[i].ID)
			transactions = append(transactions, accTransactions...)
		}

		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = transactions
		return response
	} else {
		return map[string]interface{}{"message": "Not valid token"}
	}
}
