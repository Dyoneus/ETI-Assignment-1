// By Ong Jia Yuan S10227735B
// trip-service/main.go
package main

import (
	"log"
	"net/http"
	"trip-service/database"
	"trip-service/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Retrieve the generic database object sql.DB to close it later
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get database: %v", err)
	}

	// Defer the closure of the sqlDB
	defer sqlDB.Close()

	r := mux.NewRouter()

	// Routes and handlers here
	r.HandleFunc("/trips", handlers.PublishTrip(db)).Methods("POST")
	r.HandleFunc("/trips", handlers.ListTrips(db)).Methods("GET")
	r.HandleFunc("/trips/{id:[0-9]+}", handlers.EditTrip(db)).Methods("PATCH") // Added trip ID parameter
	r.HandleFunc("/trips/{id:[0-9]+}", handlers.DeleteTrip(db)).Methods("DELETE")

	// Routes for other endpoints

	log.Println("Starting trip service on port 5001...")
	if err := http.ListenAndServe(":5001", r); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
