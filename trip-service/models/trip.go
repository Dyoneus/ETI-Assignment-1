// /trip-service/models/trip.go
package models

import (
	"time"

	"gorm.io/gorm"
)

// Trip represents a car-pooling trip published by a car owner.
type Trip struct {
	gorm.Model
	CarOwnerID         uint      `json:"car_owner_id"`
	PickUpLocation     string    `json:"pick_up_location"`
	AlternativePickUp  string    `json:"alternative_pick_up"`
	TravelStartTime    time.Time `json:"travel_start_time"`
	DestinationAddress string    `json:"destination_address"`
	AvailableSeats     int       `json:"available_seats"`
	EnrolledPassengers int       `json:"enrolled_passengers"`
}

// Reservation represents a seat reservation for a trip by a passenger.
type Reservation struct {
	gorm.Model
	TripID      uint `json:"trip_id"`
	PassengerID uint `json:"passenger_id"`
}
