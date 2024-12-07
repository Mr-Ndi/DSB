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
// PostSuggestion handles the creation of a new suggestion
func PostSuggestion(c *gin.Context) {
	var input struct {
		Content string   `json:"content" binding:"required"`
		Tags    []string `json:"tags"`
	}

	// Get claims from context
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
		return
	}

	// Extract username from claims
	claimsMap := claims.(jwt.MapClaims)
	username, ok := claimsMap["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	suggestion := models.Suggestion{
		Content:   input.Content,
		By:        username,
		Tags:      input.Tags,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	user, err := services.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if user.SuggestionCount >= 5 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You have reached the maximum number of suggestions allowed (5)"})
		return
	}

	createdSuggestion, err := services.CreateSuggestion(&suggestion)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.SuggestionCount++
	err = services.UpdateUserSuggestionCount(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update suggestion count"})
		return
	}

	c.JSON(http.StatusCreated, createdSuggestion)
}

// GetUserSuggestions godoc
// @Summary Get user's suggestion count and remaining allowed suggestions
// @Description Returns the number of suggestions made by the user and how many they can still make
// @Tags suggestions
// @Param Authorization header string true "Bearer JWT Token"
// @Success 200 {object} map[string]int "User's suggestion count and remaining allowed suggestions"
// @Failure 401 {object} map[string]string "Unauthorized - Invalid token"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /suggestions/count [get]
func GetUserSuggestions(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User unauthorized"})
		return
	}

	claimsMap := claims.(jwt.MapClaims)
	username, ok := claimsMap["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	user, err := services.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	remaining := 5 - user.SuggestionCount

	c.JSON(http.StatusOK, gin.H{
		"suggestion_count": user.SuggestionCount,
		"remaining":        remaining,
	})
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

// GetAdminSuggestions retrieves suggestions tagged for admins with status "submitted".
// @Summary Retrieve suggestions for admin
// @Description Fetch all suggestions with status "submitted" that are tagged for the logged-in admin.
// @Tags admins
// @Accept json
// @Produce json
// @Success 200 {array} models.Suggestion "List of suggestions tagged for the admin"
// @Failure 403 {object} map[string]string "Forbidden - User is not an admin"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/suggestions [get]
func GetAdminSuggestions(c *gin.Context) {
	// Extract claims from the context (assuming you have middleware that sets this up)
	userClaims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Perform a type assertion to convert userClaims to a map[string]interface{}
	claimsMap, ok := userClaims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
		return
	}

	// Check if user is an admin
	role, exists := claimsMap["role"].(string)
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not an admin"})
		return
	}

	// Define the status you want to filter by (e.g., "submitted")
	status := "submitted"

	// Fetch suggestions tagged for this admin with the specified status
	suggestions, err := services.FindSuggestionsByTag(role, status) // Adjust service function accordingly
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// RespondToSuggestion allows an admin to respond to a suggestion.
// @Summary Respond to a suggestion
// @Description Allows an admin to submit a response to a specific suggestion.
// @Tags admins
// @Accept json
// @Produce json
// @Param suggestionId path string true "Suggestion ID" // The ID of the suggestion being responded to
// @Param requestBody body models.Response true "Response content"
// @Success 201 {object} models.Response "Successfully created response"
// @Failure 403 {object} map[string]string "Forbidden - User is not an admin"
// @Failure 400 {object} map[string]string "Bad Request - Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/suggestions/{suggestionId}/respond [post]
func RespondToSuggestion(c *gin.Context) {
	// Extract claims from context (assuming you have middleware that sets this up)
	userClaims, exists := c.Get("userClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Perform type assertion on userClaims
	claimsMap, ok := userClaims.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claims format"})
		return
	}

	// Check if user is an admin
	role, exists := claimsMap["role"].(string)
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not an admin"})
		return
	}

	// Get the suggestion ID from URL parameters
	suggestionID := c.Param("suggestionId")

	// Parse the request body to get response content
	var response models.Response
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set additional fields for the response
	response.SuggestionId = suggestionID
	response.Admin = claimsMap["username"].(string) // Assuming username is in claims

	// Call service function to save the response
	createdResponse, err := services.AddResponse(&response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdResponse)
}
