package controllers

import (
	"github.com/gorilla/mux"
	"net/http"
	u "qt-api/utils"
	"qt-api/models"
	"encoding/json"
	"strconv"
)

var CreateTest = func(w http.ResponseWriter, r *http.Request) {

	test := &models.Test{}
	err := json.NewDecoder(r.Body).Decode(test)
	if err != nil {
		u.Error(w, 400, "Could not parse data")
		return
	}

	_, err = test.Create()
	if err != nil {
		u.Error(w, 500, err.Error())
	}

	u.Respond(w, 200, map[string]interface{} {"test": test})
}

var GetTest = func(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("userId") . (uint)

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	test := &models.Test{}

	if err != nil {
		u.Error(w, 400, "Bad ID in parameters")
		return
	}
	
	_, err = test.FindByID(id)

	if err != nil {
		u.Error(w, 500, err.Error())
		return
	}

	if test.UserID != userId {
		u.Error(w, 403, "Test does not belong to user")
		return
	}
	u.Respond(w, 200, map[string]interface{}{"test": test})

}

var GetTests = func(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("userId") . (uint)

	tests, err := models.FindTestsFromUser(userId)

	if err != nil {
		u.Error(w, 500, err.Error())
	}
	u.Respond(w, 200, map[string]interface{}{"tests": tests})

}