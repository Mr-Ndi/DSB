package services

import (
	"context"
	"dsb/config"
	"dsb/models"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateSuggestion creates a new suggestion with valid tags only, validated against admin roles
func CreateSuggestion(suggestion *models.Suggestion) (*models.Suggestion, error) {
	// Access the collections
	usersCollection := config.DatabaseClient.Database("DSBox").Collection("user")
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")
	adminCollection := config.DatabaseClient.Database("DSBox").Collection("admin") // Assuming this is where admin roles are stored

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate if the 'By' field references an existing user by username
	var userExists struct {
		Username string `bson:"username"`
	}
	if err := usersCollection.FindOne(ctx, bson.M{"username": suggestion.By}).Decode(&userExists); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with username %s does not exist", suggestion.By)
		}
		return nil, fmt.Errorf("error validating user: %v", err)
	}

	// Filter valid tags from the admin roles
	var validTags []string
	for _, tag := range suggestion.Tags {
		var adminExists struct {
			Role string `bson:"role"`
		}
		// Check if the tag exists as an admin role
		err := adminCollection.FindOne(ctx, bson.M{"role": tag}).Decode(&adminExists)
		if err == nil && adminExists.Role != "" {
			validTags = append(validTags, tag) // Add valid tag to the slice
		}
	}

	// Assign valid tags to the suggestion
	suggestion.Tags = validTags

	// Insert the suggestion into the database
	result, err := suggestionsCollection.InsertOne(ctx, suggestion)
	if err != nil {
		return nil, fmt.Errorf("error inserting suggestion: %v", err)
	}

	// Assign the generated ID to the suggestion
	suggestion.Id = result.InsertedID.(primitive.ObjectID)

	return suggestion, nil
}

func FindAllSuggestions() ([]models.Suggestion, error) {
	// Access the suggestions collection
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query all suggestions
	cursor, err := suggestionsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error fetching suggestions: %v", err)
	}
	defer cursor.Close(ctx)

	// Slice to hold all suggestions
	var suggestions []models.Suggestion

	// Iterate through the cursor and decode each suggestion into the slice
	for cursor.Next(ctx) {
		var suggestion models.Suggestion
		if err := cursor.Decode(&suggestion); err != nil {
			return nil, fmt.Errorf("error decoding suggestion: %v", err)
		}
		suggestions = append(suggestions, suggestion)
	}

	// Check for cursor iteration error
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %v", err)
	}

	// Return the list of suggestions
	return suggestions, nil
}

func FindSuggestionsByUser(username string) ([]models.Suggestion, error) {
	// Get the suggestions collection from the database
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")

	// Define a filter to find suggestions by username
	filter := bson.D{{Key: "by", Value: username}}

	// Options to sort by creation time, most recent first
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	// Create a context with a timeout to prevent long-running queries
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all matching suggestions
	cursor, err := suggestionsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var suggestions []models.Suggestion
	if err := cursor.All(ctx, &suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}

func FindSuggestionsWithTag(tag string) ([]models.Suggestion, error) {
	// Get the suggestions collection from the database
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")

	// Define a filter to find suggestions by tag
	filter := bson.M{"tags": tag}

	// Options to sort by creation time, most recent first
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all matching suggestions
	cursor, err := suggestionsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode all matching suggestions into a slice
	var suggestions []models.Suggestion
	if err := cursor.All(ctx, &suggestions); err != nil {
		return nil, err
	}

	return suggestions, nil
}

func AddVote(vote *models.Vote) (*models.Vote, error) {
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the necessary collections
	usersCollection := config.DatabaseClient.Database("DSBox").Collection("user")
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")
	votesCollection := config.DatabaseClient.Database("DSBox").Collection("vote")

	// Validate user (ensure 'By' is a valid username)
	var user bson.M
	if err := usersCollection.FindOne(ctx, bson.M{"username": vote.By}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid username")
		}
		return nil, fmt.Errorf("error while validating username: %v", err)
	}

	// Validate suggestion (ensure 'SuggestionId' is a valid suggestion ID)
	suggestionID, err := primitive.ObjectIDFromHex(vote.SuggestionId)
	if err != nil {
		return nil, errors.New("invalid suggestion ID format")
	}

	var suggestion bson.M
	if err := suggestionsCollection.FindOne(ctx, bson.M{"_id": suggestionID}).Decode(&suggestion); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("suggestion not found")
		}
		return nil, fmt.Errorf("error while validating suggestion: %v", err)
	}

	// Check if the user has already voted on this suggestion
	var existingVote models.Vote
	err = votesCollection.FindOne(ctx, bson.M{"by": vote.By, "suggestionId": vote.SuggestionId}).Decode(&existingVote)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("error while checking existing vote: %v", err)
	}

	// If the user has voted, check if the vote type is different
	if err == nil && existingVote.Type != vote.Type {
		// Update the existing vote with the new type
		update := bson.M{"$set": bson.M{"type": vote.Type}}
		_, err = votesCollection.UpdateOne(ctx, bson.M{"_id": existingVote.Id}, update)
		if err != nil {
			return nil, fmt.Errorf("error while updating vote: %v", err)
		}
		// Return the updated vote
		existingVote.Type = vote.Type
		return &existingVote, nil
	}

	// If the user has already voted with the same type, return an error
	if err == nil && existingVote.Type == vote.Type {
		return nil, errors.New("already voted with the same type")
	}

	// If no existing vote or same vote type, insert the new vote
	vote.Id = primitive.NewObjectID() // Assign a new ID to the vote
	_, err = votesCollection.InsertOne(ctx, vote)
	if err != nil {
		return nil, fmt.Errorf("error while inserting vote: %v", err)
	}

	// Check the status of the suggestion
	if suggestion["status"] == "pending" {
		// Count the number of upvotes for this suggestion
		upvoteCount, err := votesCollection.CountDocuments(ctx, bson.M{
			"suggestionId": vote.SuggestionId,
			"type":         "upvote",
		})
		if err != nil {
			return nil, fmt.Errorf("error while counting upvotes: %v", err)
		}

		// If the upvotes are 100 or more, update the status to "submitted"
		if upvoteCount >= 100 {
			_, err := suggestionsCollection.UpdateOne(ctx, bson.M{"_id": suggestionID}, bson.M{
				"$set": bson.M{"status": "submitted"},
			})
			if err != nil {
				return nil, fmt.Errorf("error while updating suggestion status: %v", err)
			}
		}
	}

	return vote, nil
}

// FindSuggestionsByTag retrieves suggestions tagged for a specific user and with a specific status.
func FindSuggestionsByTag(tag string, status string) ([]models.Suggestion, error) {
	var suggestions []models.Suggestion

	// Create a filter to find suggestions by tag and status
	filter := bson.M{
		"tags":   tag,
		"status": status,
	}

	// Query the database (assuming you have a MongoDB collection named "suggestions")
	cursor, err := config.DatabaseClient.Database("DSBox").Collection("suggestions").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Decode the results into the suggestions slice
	for cursor.Next(context.TODO()) {
		var suggestion models.Suggestion
		if err := cursor.Decode(&suggestion); err != nil {
			return nil, err
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

// AddResponse saves a new response to the database.
func AddResponse(response *models.Response) (*models.Response, error) {
	ctx := context.TODO()

	// Assign a new ID to the response
	response.ID = primitive.NewObjectID().Hex() // Correctly call Hex() on the ObjectID
	response.CreatedAt = time.Now()             // Set created time

	// Insert the response into the database
	_, err := config.DatabaseClient.Database("DSBox").Collection("responses").InsertOne(ctx, response)
	if err != nil {
		return nil, fmt.Errorf("error while inserting response: %v", err)
	}

	return response, nil // Return created response or handle as needed.
}

func CreateComment(by string, content string, parentID string) (*models.Suggestion, error) {
	// Access the suggestions collection
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate parent ID format
	parentObjectID, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, fmt.Errorf("invalid parent ID format: %v", err)
	}

	// Check if the parent suggestion or comment exists
	var parent models.Suggestion
	if err := suggestionsCollection.FindOne(ctx, bson.M{"_id": parentObjectID}).Decode(&parent); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("parent suggestion or comment with ID %s does not exist", parentID)
		}
		return nil, fmt.Errorf("error fetching parent suggestion: %v", err)
	}

	// Create the comment object
	comment := models.Suggestion{
		By:        by,
		Content:   content,
		Parent:    parentID,
		Tags:      nil,           // No tags for comments by default
		Reply:     "",            // No reply for comments by default
		Views:     0,             // Not applicable for comments
		Status:    parent.Status, // Inherit status from parent
		CreatedAt: time.Now(),
	}

	// Insert the comment into the collection
	result, err := suggestionsCollection.InsertOne(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("error inserting comment: %v", err)
	}

	// Assign the generated ID to the comment
	comment.Id = result.InsertedID.(primitive.ObjectID)

	return &comment, nil
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*models.User, error) {
	ctx := context.TODO() // Create a context for the database operation

	var user models.User
	// Use config.DatabaseClient instead of client
	err := config.DatabaseClient.Database("DSBox").Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user: %v", err) // Provide more context in errors
	}

	return &user, nil
}

// UpdateUserSuggestionCount updates the user's suggestion count in the database
func UpdateUserSuggestionCount(user *models.User) error {
	ctx := context.TODO() // Create a context for the database operation

	if _, err := config.DatabaseClient.Database("DSBox").Collection("users").UpdateOne(
		ctx,
		bson.M{"username": user.Username},
		bson.M{"$set": bson.M{"suggestion_count": user.SuggestionCount}},
	); err != nil {
		return fmt.Errorf("error updating user suggestion count: %v", err)
	}

	return nil
}
