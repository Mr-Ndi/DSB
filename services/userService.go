package services

import (
	"context"
	"log"
	"time"

	"dsb/config"
	"dsb/models"
)

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
