package controllers

import (
	"dsb/models"
	"dsb/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser handles user registration.
func CreateUser(c *gin.Context) {
	var user models.User
	// Bind incoming JSON to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hashing user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hashedPassword)

	// Add user to database (via service)
	inserted, err := services.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response struct excluding sensitive information
	response := struct {
		Username  string `json:"username"`
		Regnumber int    `json:"regnumber"`
	}{
		Username:  inserted.Username,
		Regnumber: inserted.Regnumber,
	}

	// Respond with a safe subset of the user data
	c.JSON(http.StatusCreated, response)
}

// HandleLogin handles user login.
func HandleLogin(c *gin.Context) {
	// Define a struct for user credentials
	type Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var credentials Credentials
	// Bind incoming JSON to credentials struct
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to validate credentials
	user, err := services.Login(credentials.Username, credentials.Password)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err.Error() == "invalid password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Handle unexpected errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send back the user information
	c.JSON(http.StatusOK, user)
}
