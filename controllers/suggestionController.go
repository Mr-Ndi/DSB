package controllers

import (
	"net/http"
	"time"

	"dsb/models"
	"dsb/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// PostSuggestion godoc
// @Summary Create a new suggestion
// @Description Add a new suggestion to the system with optional tags. Requires authentication.
// @Tags Suggestions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param suggestion body PostSuggestionInput true "Suggestion input"
// @Success 201 {object} models.Suggestion "Created suggestion"
// @Failure 400 {object} gin.H "Bad request"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /suggestion [post]
func PostSuggestion(c *gin.Context) {
	// Struct for capturing input data
	var input struct {
		Suggestion string   `json:"suggestion" binding:"required"`
		Tags       []string `json:"tags"`
	}

	// Get the claims from the context
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
		return
	}

	// Extract the username from the claims
	claimsMap := claims.(jwt.MapClaims)
	username, ok := claimsMap["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Bind the incoming JSON request body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		// Return bad request if validation fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize the Suggestion model with input fields and default values
	suggestion := models.Suggestion{
		Suggestion: input.Suggestion,
		By:         username,   // Use username instead of regnumber
		Tags:       input.Tags, // Optional field, will be empty string if not provided
		Reply:      "",         // Default
		Votes:      0,          // Default
		Views:      0,          // Default
		Status:     "pending",  // Default
		CreatedAt:  time.Now(),
	}

	// Call the service layer to save the suggestion
	createdSuggestion, err := services.CreateSuggestion(&suggestion)
	if err != nil {
		// Return internal server error if creation fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created suggestion with a 201 status code
	c.JSON(http.StatusCreated, createdSuggestion)
}

func GetAllSuggestions(c *gin.Context) {
	// Call the service to get all suggestions
	suggestions, err := services.FindAllSuggestions()
	if err != nil {
		// Return internal server error if fetching suggestions fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of suggestions with a 200 status code
	c.JSON(http.StatusOK, suggestions)
}

func GetSuggestionsByUser(c *gin.Context) {
	// Extract the username from the path parameters
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Call the service to get suggestions by the username
	suggestions, err := services.FindSuggestionsByUser(username)
	if err != nil {
		// Return internal server error if fetching suggestions fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of suggestions made by the user
	c.JSON(http.StatusOK, suggestions)
}

// PostSuggestionInput is the struct used to document the input for the PostSuggestion endpoint
type PostSuggestionInput struct {
	Suggestion string   `json:"suggestion" example:"This is a suggestion"` // Example input
	By         int      `json:"by" example:"123"`                          // Example input
	Tags       []string `json:"tags" example:"[\"tag1\", \"tag2\"]"`       // Example input
}
