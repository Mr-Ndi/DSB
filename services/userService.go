package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"dsb/config"
	"dsb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(username, password string) (*models.User, error) {
	collection := config.DatabaseClient.Database("DSBox").Collection("user")

	filter := bson.M{"username": username} // Filter to search by username

	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Handle the case where no user is found
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("Error querying database: %v", err)
		return nil, err
	}

	// Log the user (excluding sensitive info)
	fmt.Printf("User found: %+v\n", user)

	// Compare the hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func AddUser(user *models.User) (*models.User, error) {

	if config.DatabaseClient == nil {
		log.Fatal("MongoDB client is nil")
	}

	collection := config.DatabaseClient.Database("DSBox").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
