package controllers

import (
	"dsb/models"
	"dsb/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Credentials represents the structure for user login credentials.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUser handles user registration.
// @Summary Create a new user
// @Description Create a new user with username, registration number and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user [post]
func CreateUser(c *gin.Context) {
	var user models.User
	// Bind incoming JSON to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hashing user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
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
// @Summary User login
// @Description Login using username and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User credentials"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func HandleLogin(c *gin.Context) {
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

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// AdminRequest represents the structure for admin registration.
type AdminRequest struct {
	Role     string `json:"role"`     // Single role for the admin
	Password string `json:"password"` // Admin password
}

// RegisterAdmin handles admin registration.
// @Summary Register a new admin
// @Description Create a new admin with a role and password
// @Tags admins
// @Accept json
// @Produce json
// @Param admin body AdminRequest true "Admin data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin [post]
// RegisterAdmin handles admin registration.
func RegisterAdmin(c *gin.Context) {
	var adminRequest AdminRequest

	// Bind JSON to AdminRequest struct
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		log.Printf("Error binding JSON: %v", err) // Log the error
		return
	}

	// Validate role
	allowedRoles := map[string]bool{
		"admin":     true,
		"moderator": true,
	}
	if !allowedRoles[adminRequest.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role provided"})
		return
	}

	// Validate password
	if len(adminRequest.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		log.Printf("Error hashing password: %v", err) // Log the error
		return
	}

	// Create admin instance
	admin := models.Admin{
		Role:     adminRequest.Role,
		Password: string(hashedPassword),
	}

	// Add admin to the database
	inserted, err := services.AddAdmin(&admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		log.Printf("Error inserting admin into database: %v", err) // Log the error
		return
	}

	// Respond with the created admin's role
	c.JSON(http.StatusCreated, gin.H{
		"role":    inserted.Role,
		"message": "Admin created successfully",
	})
}
