package handlers

import (
	"encoding/json"
	"net/http"
	"user-service/models"

	"gorm.io/gorm"
)

var db *gorm.DB // Add this line to declare the db variable

// CreateUser creates a new user in the database.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// You will need to decode the request body into the User struct.
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// insert the new User into the database.
	db.Create(&user)

	// Finally, encode the created user into JSON and send it in the response.
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
