package utils

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func Respond(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, statusCode int, statusMsg string) {
	fmt.Println(statusCode, statusMsg)
	Respond(w, statusCode, struct {
				Error string `json:"error"`
			}{
				Error: statusMsg,
			})
}