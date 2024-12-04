package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseClient *mongo.Client
var JwtKey string

// DatabaseConnection establishes a connection to the MongoDB database and returns an error if it fails.
func DatabaseConnection() error {
	// Retrieve JWT key and database URI from environment variables
	JwtKey = os.Getenv("JWT_KEY")
	dbURI := os.Getenv("DB_URI")

	if JwtKey == "" {
		return fmt.Errorf("Jwt secret key missing")
	}

	if dbURI == "" {
		return fmt.Errorf("Connection string is empty")
	}

	clientOptions := options.Client().ApplyURI(dbURI)

	timeOut, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(timeOut, clientOptions)
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}

	err = client.Ping(timeOut, nil)
	if err != nil {
		return fmt.Errorf("Failed to ping connection to db: %v", err)
	}

	fmt.Println("Successfully connected to database!")
	DatabaseClient = client
	return nil // Return nil if successful
}
