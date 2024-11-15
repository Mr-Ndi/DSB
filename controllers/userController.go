// controllers/userController.go
package controllers

import (
	"context"
	"dsb/models"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles user registration.
func CreateUser(w http.ResponseWriter, r *http.Request, userCollection *mongo.Collection) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hashing user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert user into database
	_, err = userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response struct to exclude sensitive information
	response := struct {
		Username  string `json:"username"`
		Regnumber int    `json:"regnumber"`
	}{
		Username:  user.Username,
		Regnumber: user.Regnumber,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response) // Encode and send only safe fields
}
