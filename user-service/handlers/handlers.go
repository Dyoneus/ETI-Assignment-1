// /user-service/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
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
