package controllers

import (
	"fmt"
	"net/http"
	"time"

	"dsb/models"
	"dsb/services"

	"github.com/gin-gonic/gin"
)

// postSuggestion handles the request to create a new suggestion
func PostSuggestion(c *gin.Context) {
	// Struct for capturing input data
	var input struct {
		Suggestion string   `json:"suggestion" binding:"required"`
		By         int      `json:"by" binding:"required"`
		Tags       []string `json:"tags"`
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
		By:         input.By,
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
