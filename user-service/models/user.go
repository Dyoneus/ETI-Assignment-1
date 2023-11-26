// /user-service/models/user.go
package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system with a default passenger profile.
type User struct {
	gorm.Model
	DeletedAt gorm.DeletedAt `gorm:"index"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Mobile    string         `json:"mobile"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Password  string         `json:"password"`
	UserType  string         `json:"user_type"`
	Active    bool           `gorm:"default:true"`
}

// CarOwnerProfile represents a car owner's profile, extending the User model.
type CarOwnerProfile struct {
	UserID         uint   `gorm:"primaryKey;autoIncrement:false"`
	DriversLicense string `json:"drivers_license"`
	CarPlateNumber string `json:"car_plate_number"`
}
