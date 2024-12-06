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
// @Description Add a new suggestion for the authenticated user
// @Tags suggestions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT Token"
// @Param suggestion body PostSuggestionInput true "Suggestion Details"
// @Success 201 {object} models.Suggestion "Successfully created suggestion"
// @Failure 400 {object} map[string]string "Bad Request - Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized - Invalid token"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /suggestions [post]
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

// GetAllSuggestions godoc
// @Summary Retrieve all suggestions
// @Description Get a list of all suggestions in the system
// @Tags suggestions
// @Produce json
// @Success 200 {array} models.Suggestion "Successfully retrieved suggestions"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /suggestions [get]
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

// GetSuggestionsByUser godoc
// @Summary Get suggestions by username
// @Description Retrieve all suggestions submitted by a specific user
// @Tags suggestions
// @Produce json
// @Param username path string true "Username"
// @Success 200 {array} models.Suggestion "Successfully retrieved user suggestions"
// @Failure 400 {object} map[string]string "Bad Request - Missing username"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /suggestions/user/{username} [get]
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

// GetSuggestionWithTag godoc
// @Summary Get suggestions by tag/role
// @Description Retrieve suggestions filtered by a specific tag or role
// @Tags suggestions
// @Produce json
// @Param role path string true "Role/Tag"
// @Success 200 {array} models.Suggestion "Successfully retrieved suggestions"
// @Success 404 {object} map[string]string "No suggestions found for the tag"
// @Failure 400 {object} map[string]string "Bad Request - Missing role"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /suggestions/tag/{role} [get]
func GetSuggestionWithTag(c *gin.Context) {
	// Get the tag (role) from path parameters
	tag := c.Param("role")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}

	// Call the service to find suggestions with the given tag
	suggestions, err := services.FindSuggestionsWithTag(tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suggestions"})
		return
	}

	// If no suggestions are found, return an appropriate response
	if len(suggestions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No suggestions found for the specified role"})
		return
	}

	// Return the list of suggestions with a 200 status code
	c.JSON(http.StatusOK, suggestions)
}

// HandleVote godoc
// @Summary Add or update a vote on a suggestion
// @Description Allows a user to upvote or downvote a suggestion. If the user has already voted, it will update the vote type.
// @Tags votes
// @Accept json
// @Produce json
// @Param type path string true "Vote type (upvote or downvote)"
// @Param vote body object true "Vote data with suggestion ID" example({"suggestionId": "12345"})
// @Success 200 {object} map[string]string "Successfully added or updated vote"
// @Failure 400 {object} map[string]string "Invalid request body or vote type"
// @Failure 404 {object} map[string]string "Username missing in claims, or suggestion not found"
// @Failure 409 {object} map[string]string "You have already voted with the same type"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /votes/{type} [post]
//
// @BodyExample json { "suggestionId": "12345" }
func HandleVote(c *gin.Context) {
	// Retrieve the claims (which contain the username) from the context
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Type assert the claims to the custom struct that contains the username
	userClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Extract the username from the claims
	username, exists := userClaims["username"].(string)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username missing in claims"})
		return
	}

	// Get the vote type from the URL parameters
	voteType := c.Param("type")
	if voteType != "upvote" && voteType != "downvote" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote type"})
		return
	}

	// Parse the request body to get the suggestion ID
	var requestBody struct {
		SuggestionId string `json:"suggestionId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create the vote object
	vote := &models.Vote{
		By:           username,
		Type:         voteType,
		SuggestionId: requestBody.SuggestionId,
	}

	// Call the AddVote function to process the vote
	createdVote, err := services.AddVote(vote)
	if err != nil {
		// Handle specific errors based on the returned error message
		if err.Error() == "invalid username" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Username not found"})
		} else if err.Error() == "invalid suggestion ID format" || err.Error() == "suggestion not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Suggestion not found"})
		} else if err.Error() == "already voted with the same type" {
			c.JSON(http.StatusConflict, gin.H{"error": "You have already voted with the same type"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while processing your vote"})
		}
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, createdVote)
}

// PostSuggestionInput is the struct used to document the input for the PostSuggestion endpoint
type PostSuggestionInput struct {
	Suggestion string   `json:"suggestion" example:"This is a suggestion"` // Example input
	By         int      `json:"by" example:"123"`                          // Example input
	Tags       []string `json:"tags" example:"[\"tag1\", \"tag2\"]"`       // Example input
}
