// By Ong Jia Yuan S10227735B
// main.go in user-service directory
// user-service/main.go

package main

import (
	"log"
	"net/http"
	"user-service/database" // This imports the database package where InitializeDatabase is defined
	"user-service/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Use the InitializeDatabase function from the database package to set up the database connection.
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get database: %v", err)
	}

	defer sqlDB.Close()

	// Set up the router.
	r := mux.NewRouter()

	// Handlers
	r.HandleFunc("/login", handlers.Login(db)).Methods("POST")
	r.HandleFunc("/users", handlers.CreateUser(db)).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserByID(db)).Methods("GET")

	// Handlers for updating profile
	r.HandleFunc("/users", handlers.UpdateUser(db)).Methods("PATCH")

	// Start the server.
	log.Println("Starting user service on port 5000...")
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
