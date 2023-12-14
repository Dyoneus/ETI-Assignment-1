// /trip-service/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"trip-service/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func PublishTrip(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trip models.Trip
		if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

// ListTrips handler
func ListTrips(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request method and URL
		//fmt.Printf("Received %s request to %s\n", r.Method, r.URL.String())

		carOwnerID := r.URL.Query().Get("carOwnerID") // Get the carOwnerID from the query parameters

		var trips []models.Trip
		if result := db.Where("car_owner_id = ?", carOwnerID).Find(&trips); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trips)
	}
}

// EditTrip handler
func EditTrip(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trip models.Trip
		if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Perform update operation, assuming trip.ID is set
		if result := db.Save(&trip); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Trip updated successfully")
	}
}

// DeleteTrip handler
func DeleteTrip(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tripID := vars["id"]

		// Perform delete operation
		if result := db.Delete(&models.Trip{}, tripID); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Trip deleted successfully")
	}
}

func ListSoftDeletedTrips(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trips []models.Trip
		if result := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&trips); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trips)
	}
}

func AvailableTrips(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trips []models.Trip
		// Assuming there's a DeletedAt field in the model to mark soft deleted records
		// and an AvailableSeats field to show how many seats are left.
		result := db.Where("deleted_at IS NULL AND available_seats > 0").Find(&trips)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(trips); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
