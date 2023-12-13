// /trip-service/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"trip-service/models"

	"gorm.io/gorm"
)

func PublishTrip(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trip models.Trip
		if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Add validation logic for trip fields...

		// Save the trip to the database
		if result := db.Create(&trip); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(trip)
	}
}
