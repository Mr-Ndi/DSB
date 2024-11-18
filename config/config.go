package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseClient *mongo.Client

func DatabaseConnection() {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file:%v", err)
	}

	dbURI := os.Getenv("DB_URI")

	if dbURI == "" {
		log.Fatalf("Connection string is empty")
	}

	clientOptions := options.Client().ApplyURI(dbURI)

	timeOut, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(timeOut, clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = client.Ping(timeOut, nil)

	if err != nil {
		log.Fatalf("Failed to ping connection to db: %v", err)
	}

	fmt.Println("Successfully connected to database !!")

	DatabaseClient = client
}
