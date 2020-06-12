package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/loxt/go-fintech-banking/helpers"
	"github.com/loxt/go-fintech-banking/users"
	"io/ioutil"
	"log"
	"net/http"
)

type Login struct {
	Username string
	Password string
}

type ErrResponse struct {
	Message string
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)

	if login["message"] == "All is fine" {
		res := login
		json.NewEncoder(w).Encode(res)
	} else {
		res := ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(res)
	}
}

func StartApi() {
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	fmt.Println("App is working on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
