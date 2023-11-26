// /user-service/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"
	"user-service/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt" // Hashing password
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Server-side validation
		if user.FirstName == "" || user.LastName == "" || user.Mobile == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "All fields are required and cannot be empty.", http.StatusBadRequest)
			return
		}

		// Set default user type to "passenger"
		user.UserType = "passenger"

		// Hash the user's password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)

		// Create the user in the database
		if result := db.Create(&user); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Encode and send the created user as the response.
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define a struct to decode the login request
		var loginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Find the user by email
		var user models.User
		result := db.Where("email = ?", loginRequest.Email).First(&user)
		if result.Error != nil {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		// Compare the hashed password
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if err != nil {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		// If login is successful, create a response struct with the UserType
		response := struct {
			UserType string `json:"userType"`
		}{
			UserType: user.UserType,
		}

		// Set Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send the response with the UserType
		json.NewEncoder(w).Encode(response)
	}
}

// getUsers returns a handler function that retrieves all users from the database.
func GetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		result := db.Find(&users)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserByID(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var user models.User
		result := db.First(&user, id)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateRequest struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updateRequest.Email == "" {
			http.Error(w, "Email address is required", http.StatusBadRequest)
			return
		}

		var user models.User
		result := db.Where("email = ? AND deleted_at IS NULL", updateRequest.Email).First(&user)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		user.FirstName = updateRequest.FirstName
		user.LastName = updateRequest.LastName

		if err := db.Save(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUserMobile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define a struct to decode the request body
		var updateMobileRequest struct {
			Email  string `json:"email"`
			Mobile string `json:"mobile"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateMobileRequest); err != nil {
			http.Error(w, "Error decoding request: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Perform the update operation
		result := db.Model(&models.User{}).
			Where("email = ?", updateMobileRequest.Email).
			Update("mobile", updateMobileRequest.Mobile)
		if result.Error != nil {
			http.Error(w, "Failed to update mobile number: "+result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Send a response back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"result": "Mobile number updated successfully",
		})
	}
}

// isValidEmail checks if the provided email address has a valid format.
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func UpdateUserEmail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var emailUpdateRequest struct {
			OldEmail string `json:"old_email"`
			NewEmail string `json:"new_email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&emailUpdateRequest); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate the new email
		if !isValidEmail(emailUpdateRequest.NewEmail) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Update the user's email
		result := db.Model(&models.User{}).
			Where("email = ?", emailUpdateRequest.OldEmail).
			Update("email", emailUpdateRequest.NewEmail)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Email updated successfully")
	}
}

// UpdateDriversLicense updates the driver's license number of the car owner profile.
func UpdateDriversLicense(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body
		var updateRequest struct {
			Email          string `json:"email"`
			DriversLicense string `json:"drivers_license"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the input
		if updateRequest.Email == "" || updateRequest.DriversLicense == "" {
			http.Error(w, "Email and driver's license number are required", http.StatusBadRequest)
			return
		}

		// Update the driver's license number in the car_owner_profiles table
		result := db.Model(&models.CarOwnerProfile{}).
			Where("user_id = (SELECT id FROM users WHERE email = ?)", updateRequest.Email).
			Update("drivers_license", updateRequest.DriversLicense)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			http.Error(w, "No car owner profile found for the given email", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Driver's license updated successfully")
	}
}

// UpdateCarPlate updates the car plate number of the car owner profile.
func UpdateCarPlate(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body
		var updateRequest struct {
			Email          string `json:"email"`
			CarPlateNumber string `json:"car_plate_number"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate the input
		if updateRequest.Email == "" || updateRequest.CarPlateNumber == "" {
			http.Error(w, "Email and car plate number are required", http.StatusBadRequest)
			return
		}

		// Update the car plate number in the car_owner_profiles table
		result := db.Model(&models.CarOwnerProfile{}).
			Where("user_id = (SELECT id FROM users WHERE email = ?)", updateRequest.Email).
			Update("car_plate_number", updateRequest.CarPlateNumber)

		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			http.Error(w, "No car owner profile found for the given email", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Car plate number updated successfully")
	}
}

func DeleteUserAccount(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Assume we have a field `Active` in the User model to mark the account as inactive
		user.Active = false
		now := time.Now()
		user.DeletedAt = gorm.DeletedAt{Time: now, Valid: true} // Set the DeletedAt to current time for soft delete

		if err := db.Save(&user).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "User account deleted successfully")
	}
}

func UpgradeToCarOwner(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgradeRequest struct {
			Email          string `json:"email"`
			DriversLicense string `json:"drivers_license"`
			CarPlateNumber string `json:"car_plate_number"`
		}

		if err := json.NewDecoder(r.Body).Decode(&upgradeRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Retrieve the userID from the user's email
		var user models.User
		if err := db.Where("email = ?", upgradeRequest.Email).First(&user).Error; err != nil {
			// Handle error, e.g., user not found
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Begin a transaction
		tx := db.Begin()

		// Insert into car_owner_profiles
		carOwnerProfile := models.CarOwnerProfile{
			UserID:         user.ID, // Use the ID from the User struct
			DriversLicense: upgradeRequest.DriversLicense,
			CarPlateNumber: upgradeRequest.CarPlateNumber,
		}

		if err := tx.Create(&carOwnerProfile).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to insert car owner profile: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Update the user_type in users table
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("user_type", "car_owner").Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to update user type to car owner: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			http.Error(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "User upgraded to car owner successfully")
	}
}
