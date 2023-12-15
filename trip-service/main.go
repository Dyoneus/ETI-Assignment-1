// By Ong Jia Yuan S10227735B
// trip-service/main.go
package main

import (
	"log"
	"net/http"
	"time"
	"trip-service/database"
	"trip-service/handlers"
	"trip-service/models"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Check current time to published trip that is past scheduled time
func scheduleTripDeletion(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			var trips []models.Trip
			db.Where("travel_start_time < ? AND deleted_at IS NULL", now).Find(&trips)
			for _, trip := range trips {
				db.Model(&trip).Update("DeletedAt", now)
			}
		}
	}
}

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

	// Start the routine to schedule trip deletions
	go scheduleTripDeletion(db)

	// Setup CORS
	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := gorillaHandlers.AllowedOrigins([]string{"*"})
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "PATCH", "POST", "DELETE", "OPTIONS"})

	corsHandler := gorillaHandlers.CORS(originsOk, headersOk, methodsOk)

	r := mux.NewRouter()

	// Routes and handlers here
	r.HandleFunc("/trips", handlers.PublishTrip(db)).Methods("POST")
	r.HandleFunc("/trips", handlers.ListTrips(db)).Methods("GET")
	r.HandleFunc("/trips/{id:[0-9]+}", handlers.EditTrip(db)).Methods("PATCH") // Added trip ID parameter
	r.HandleFunc("/trips/{id:[0-9]+}", handlers.GetTrip(db)).Methods("GET")
	r.HandleFunc("/trips/{id:[0-9]+}", handlers.DeleteTrip(db)).Methods("DELETE")
	r.HandleFunc("/past-trips", handlers.ListSoftDeletedTrips(db)).Methods("GET")
	r.HandleFunc("/available-trips", handlers.AvailableTrips(db)).Methods("GET")
	r.HandleFunc("/enroll", handlers.EnrollInTrip(db)).Methods("POST")
	r.HandleFunc("/enrolled-trips", handlers.GetEnrolledTripsHandler(db)).Methods("GET")
	r.HandleFunc("/past-trips/passenger", handlers.GetPastTripsForPassenger(db)).Methods("GET")

	log.Println("Starting trip service on port 5001...")
	if err := http.ListenAndServe(":5001", corsHandler(r)); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
