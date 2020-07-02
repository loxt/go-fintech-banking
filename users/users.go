package users

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/loxt/go-fintech-banking/database"
	"github.com/loxt/go-fintech-banking/helpers"
	"github.com/loxt/go-fintech-banking/interfaces"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute ^ 60).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)

	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)
	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var response = map[string]interface{}{"message": "All is fine"}
	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response
}

func Login(username string, pass string) map[string]interface{} {
	valid := helpers.Validation([]interfaces.Validation{
		{
			Value: username,
			Valid: "username",
		},
		{
			Value: pass,
			Valid: "password",
		},
	})

	if valid {
		user := &interfaces.User{}

		if database.DB.Where("username = ?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}

		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		var accounts []interfaces.ResponseAccount
		database.DB.Table("accounts").
			Select("id, name, balance").
			Where("user_id = ?", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, true)

		return response
	} else {
		return map[string]interface{}{"message": "Not valid values"}
	}
}

func Register(username string, email string, pass string) map[string]interface{} {
	valid := helpers.Validation([]interfaces.Validation{
		{
			Value: username,
			Valid: "username",
		},
		{
			Value: email,
			Valid: "email",
		},
		{
			Value: pass,
			Valid: "password",
		},
	})

	if valid {
		generatePassword := helpers.HashAndSalt([]byte(pass))
		user := &interfaces.User{
			Username: username,
			Email:    email,
			Password: generatePassword,
		}
		database.DB.Create(&user)

		account := &interfaces.Account{
			Type:    "Daily Account",
			Name:    username + "'s account",
			Balance: 0,
			UserID:  user.ID,
		}

		database.DB.Create(&account)

		var accounts []interfaces.ResponseAccount
		respAccount := interfaces.ResponseAccount{
			ID:      account.ID,
			Name:    account.Name,
			Balance: int(account.Balance),
		}

		accounts = append(accounts, respAccount)

		var response = prepareResponse(user, accounts, true)
		return response
	} else {
		return map[string]interface{}{"message": "Not valid values"}
	}
}

func GetUser(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)

	if isValid {
		user := &interfaces.User{}

		if database.DB.Where("id = ?", id).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}

		var accounts []interfaces.ResponseAccount
		database.DB.Table("accounts").
			Select("id, name, balance").
			Where("user_id = ?", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, false)
		return response
	} else {
		return map[string]interface{}{"message": "Not valid token"}
	}
}
