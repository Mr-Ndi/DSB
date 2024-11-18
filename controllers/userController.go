// controllers/userController.go
package controllers

import (
	"dsb/models"
	"dsb/services"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles user registration.
func CreateUser(res http.ResponseWriter, req *http.Request) {
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Hashing user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	inserted, err := services.AddUser(&user)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response struct to exclude sensitive information
	response := struct {
		Username  string `json:"username"`
		Regnumber int    `json:"regnumber"`
	}{
		Username:  inserted.Username,
		Regnumber: inserted.Regnumber,
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(response) // Encode and send only safe fields
}
