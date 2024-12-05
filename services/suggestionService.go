package services

import (
	"context"
	"dsb/config"
	"dsb/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Check if the 'by' field references an existing user
	var userExists struct {
		Id int `bson:"regnumber"`
	}

	err := usersCollection.FindOne(ctx, bson.M{"regnumber": suggestion.By}).Decode(&userExists)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user with ID %d does not exist", suggestion.By)
		}
		return nil, err
	}

	// Filter out invalid tags
	validTags := []string{}
	for _, tag := range suggestion.Tags {
		var adminExists struct {
			Role string `bson:"role"`
		}

		// Check if the tag exists as an admin role
		err := adminCollection.FindOne(ctx, bson.M{"role": tag}).Decode(&adminExists)
		if err == nil && adminExists.Role != "" {
			validTags = append(validTags, tag) // Valid tag, add it to the validTags slice
		}
	}

	// If there are valid tags, set them; otherwise, set an empty array
	suggestion.Tags = validTags

	// Insert the suggestion into the database
	result, err := suggestionsCollection.InsertOne(ctx, suggestion)
	if err != nil {
		return nil, err
	}

	// Assign the generated ID to the suggestion
	suggestion.Id = result.InsertedID.(primitive.ObjectID)

	return suggestion, nil
}
