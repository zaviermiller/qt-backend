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
		u.Error(w, 400, "Bad Request")
		return
	}

	_, err = user.Create()
	if err != nil {
		u.Error(w, 500, err.Error())
	}
	u.Respond(w, 200, "")

}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		u.Error(w, 400, "Bad Request")
		return
	}

	token, err := models.Login(user.Email, user.Password)
	
	if err != nil {
		u.Error(w,500, err.Error())
	}

	u.Respond(w, 200, map[string]interface{}{"token": token})

}

var ConfirmUser = func(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("userId") . (uint)
	user := *models.GetUser(userId)

	u.Respond(w, 200, map[string]interface{}{"user":user})

}