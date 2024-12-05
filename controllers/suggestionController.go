package controllers

import (
	"fmt"
	"net/http"
	"time"

	"dsb/models"
	"dsb/services"

	"github.com/gin-gonic/gin"
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

	regnumber, ok := c.Get("regnumber")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
	}

	regnumberFloat, ok := regnumber.(float64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
	}

	regnumberInt := int(regnumberFloat)

	// Bind the incoming JSON request body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		// Return bad request if validation fails
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize the Suggestion model with input fields and default values
	suggestion := models.Suggestion{
		Suggestion: input.Suggestion,
		By:         regnumberInt,
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

	fmt.Printf(createdSuggestion.Suggestion)

	// Return the created suggestion with a 201 status code
	c.JSON(http.StatusCreated, createdSuggestion)
}

// PostSuggestionInput is the struct used to document the input for the PostSuggestion endpoint
type PostSuggestionInput struct {
	Suggestion string   `json:"suggestion" example:"This is a suggestion"` // Example input
	By         int      `json:"by" example:"123"`                          // Example input
	Tags       []string `json:"tags" example:"[\"tag1\", \"tag2\"]"`       // Example input
}
