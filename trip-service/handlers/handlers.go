// /trip-service/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

		// Log the trip object to check if ID is zero
		//fmt.Printf("Received trip for update: %+v\n", trip)

		/*
			// Log the trip ID and check it's not zero
			fmt.Printf("Attempting to update trip with ID: %d\n", trip.ID)
			if trip.ID == 0 {
				http.Error(w, "Trip ID is zero", http.StatusBadRequest)
				return
			}
		*/

		// Fetch the existing trip from the database
		var existingTrip models.Trip
		if err := db.First(&existingTrip, trip.ID).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Update the trip with new data
		if err := db.Model(&existingTrip).Updates(trip).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Trip updated successfully")

		/*
			// Perform update operation, assuming trip.ID is set
			if result := db.Save(&trip); result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Trip updated successfully")
		*/
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

func EnrollInTrip(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var enrollmentData struct {
			PassengerID uint `json:"passenger_id"`
			TripID      uint `json:"trip_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&enrollmentData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Start a transaction
		tx := db.Begin()

		var trip models.Trip
		if err := tx.First(&trip, enrollmentData.TripID).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		}

		// Check if there are available seats
		if trip.AvailableSeats <= 0 {
			tx.Rollback()
			http.Error(w, "No available seats", http.StatusConflict)
			return
		}

		// Check if the user is already enrolled in the trip
		if alreadyEnrolled(db, enrollmentData.PassengerID, enrollmentData.TripID) {
			http.Error(w, "User is already enrolled in this trip", http.StatusBadRequest)
			return
		}

		// Reduce available seats and increase enrolled passengers
		trip.AvailableSeats--
		trip.EnrolledPassengers++

		if err := tx.Save(&trip).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update trip", http.StatusInternalServerError)
			return
		}

		// Create a reservation
		reservation := models.Reservation{
			TripID:      trip.ID,
			PassengerID: enrollmentData.PassengerID,
		}

		if err := tx.Create(&reservation).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		tx.Commit()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Enrolled in trip successfully")
	}
}

func alreadyEnrolled(db *gorm.DB, passengerID, tripID uint) bool {
	// Query the reservations table for an entry with passengerID and tripID
	var count int64
	db.Table("reservations").Where("passenger_id = ? AND trip_id = ?", passengerID, tripID).Count(&count)
	return count > 0
}

// Get all enrolled trips for a passenger
func GetEnrolledTripsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract passengerID from the query parameters
		passengerID := r.URL.Query().Get("passenger_id") // Make sure this matches the frontend key

		if passengerID == "" {
			http.Error(w, "Passenger ID is required", http.StatusBadRequest)
			return
		}

		// Convert passengerID to uint
		pid, err := strconv.ParseUint(passengerID, 10, 32)
		if err != nil {
			http.Error(w, "Invalid Passenger ID", http.StatusBadRequest)
			return
		}

		// Find all reservations for the passenger
		var reservations []models.Reservation
		if result := db.Where("passenger_id = ?", pid).Find(&reservations); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// If no reservations found, return an empty array
		if len(reservations) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]models.Trip{})
			return
		}

		// Extract trip IDs from reservations
		var tripIDs []uint
		for _, reservation := range reservations {
			tripIDs = append(tripIDs, reservation.TripID)
		}

		// Find all trips corresponding to the trip IDs
		var trips []models.Trip
		if result := db.Where("id IN ?", tripIDs).Find(&trips); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Send the trips back to the client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trips)
		/*
			// Extract passengerID from the query parameters
			passengerID := r.URL.Query().Get("passengerID")

			// Find all reservations for the passenger
			var reservations []models.Reservation
			db.Where("passenger_id = ?", passengerID).Find(&reservations)

			// Extract trip IDs from reservations
			var tripIDs []uint
			for _, reservation := range reservations {
				tripIDs = append(tripIDs, reservation.TripID)
			}

			// Find all trips corresponding to the trip IDs
			var trips []models.Trip
			db.Where("id IN ?", tripIDs).Find(&trips)

			// Send the trips back to the client
			json.NewEncoder(w).Encode(trips)
		*/
	}
}

func GetPastTripsForPassenger(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract passengerID from query parameters
		passengerIDParam := r.URL.Query().Get("passengerID")
		if passengerIDParam == "" {
			http.Error(w, "Passenger ID is required", http.StatusBadRequest)
			return
		}

		// Convert passengerID to uint
		passengerID, err := strconv.ParseUint(passengerIDParam, 10, 32)
		if err != nil {
			http.Error(w, "Invalid Passenger ID", http.StatusBadRequest)
			return
		}

		// Query for past trips including soft-deleted ones
		var trips []models.Trip
		result := db.Unscoped().Table("trips").Joins("JOIN reservations ON reservations.trip_id = trips.id").
			Where("reservations.passenger_id = ?", passengerID).
			Find(&trips)

		// Handle possible errors from the database query
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Return the trips as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trips)
	}
}

func ViewPastTripsForPassenger(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		passengerID := r.URL.Query().Get("passengerID")
		if passengerID == "" {
			http.Error(w, "Passenger ID is required", http.StatusBadRequest)
			return
		}

		var trips []models.Trip
		// Assuming `models.Reservation` has a `PassengerID` field and a `TripID` field.
		// The `JOIN` operation is based on matching `TripID` in both `trips` and `reservations`.
		result := db.Unscoped().Joins("JOIN reservations ON reservations.trip_id = trips.id").
			Where("reservations.passenger_id = ?", passengerID).
			Find(&trips)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trips)
	}
}

// GetTrip handles GET requests for a single trip by ID
func GetTrip(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tripID := vars["id"]

		var trip models.Trip
		if result := db.First(&trip, "id = ?", tripID); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(w, "Trip not found", http.StatusNotFound)
			} else {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trip)
	}
}

func CancelEnrollment(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cancelData struct {
			PassengerID uint `json:"passenger_id"`
			TripID      uint `json:"trip_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&cancelData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Start a transaction
		tx := db.Begin()

		// Find the trip to ensure it exists
		var trip models.Trip
		if err := tx.First(&trip, cancelData.TripID).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Trip not found", http.StatusNotFound)
			return
		}

		// Check if the user is enrolled in the trip
		var reservation models.Reservation
		if err := tx.Where("passenger_id = ? AND trip_id = ?", cancelData.PassengerID, cancelData.TripID).First(&reservation).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Reservation not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Mark the reservation as deleted
		if err := tx.Model(&reservation).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to cancel enrollment", http.StatusInternalServerError)
			return
		}

		// Update the trip's available seats and enrolled passengers
		if err := tx.Model(&trip).Updates(map[string]interface{}{
			"available_seats":     gorm.Expr("available_seats + ?", 1),
			"enrolled_passengers": gorm.Expr("enrolled_passengers - ?", 1),
		}).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update trip", http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Canceled enrollment successfully")
	}
}
