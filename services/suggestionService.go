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

func CreateSuggestion(suggestion *models.Suggestion) (*models.Suggestion, error) {

	// Access the collections
	usersCollection := config.DatabaseClient.Database("DSBox").Collection("user")
	suggestionsCollection := config.DatabaseClient.Database("DSBox").Collection("suggestion")

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

	// Insert the suggestion into the database
	result, err := suggestionsCollection.InsertOne(ctx, suggestion)
	if err != nil {
		return nil, err
	}

	// Assign the generated ID to the suggestion
	suggestion.Id = result.InsertedID.(primitive.ObjectID)

	return suggestion, nil
}
