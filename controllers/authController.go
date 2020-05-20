package controllers

import (
	"net/http"
	u "qt-api/utils"
	"qt-api/models"
	"encoding/json"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	response := user.Create()
	u.Respond(w, response)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	response := models.Login(user.Email, user.Password)
	u.Respond(w, response)
}