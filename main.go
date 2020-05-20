package main

import (
	"github.com/gorilla/mux"
	"qt-api/app"
	"qt-api/controllers"
	"os"
	"fmt"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("[*] QuickTests Go API running!")

	err := http.ListenAndServe(":"  + port, router)
	if err != nil {
		fmt.Print(err)
	}

	router.HandleFunc("/v1/users/create", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/v1/users/authenticate", controllers.Authenticate).Methods("POST")
}
