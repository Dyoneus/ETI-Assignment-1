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
	CarOwnerName       string    `json:"car_owner_name"`
	PickUpLocation     string    `json:"pick_up_location"`
	AlternativePickUp  string    `json:"alternative_pick_up"`
	DestinationAddress string    `json:"destination_address"`
	TravelStartTime    time.Time `json:"travel_start_time"`
	AvailableSeats     int       `json:"available_seats"`
	EnrolledPassengers int       `json:"enrolled_passengers"`

	Reservation []Reservation `gorm:"foreignKey:TripID"`
}

// Reservation represents a seat reservation for a trip by a passenger.
type Reservation struct {
	gorm.Model
	TripID      uint `json:"trip_id"` // Foreign key
	PassengerID uint `json:"passenger_id"`

	// Rreference to the Trip struct
	Trip Trip `gorm:"references:ID"`
}
