package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"qt-api/app"
	"qt-api/controllers"
	"os"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("[*] QuickTests Go API starting...")

	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)
	
	// User routes
	router.HandleFunc("/v1/users/create", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/v1/users/authenticate", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/v1/user", controllers.ConfirmUser).Methods("GET")

	// Test routes
	router.HandleFunc("/v1/tests/create", controllers.CreateTest).Methods("POST")
	router.HandleFunc("/v1/tests/{id}", controllers.GetTest).Methods("GET")
	router.HandleFunc("/v1/tests", controllers.GetTests).Methods("GET")	
	router.HandleFunc("/v1/teststest", controllers.GetTests).Methods("GET")

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"POST", "GET", "PUT", "DELETE"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("[*] QuickTests API v1 started")

	err := http.ListenAndServe(":"  + port, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))
	if err != nil {
		fmt.Println(err)
	}

}
