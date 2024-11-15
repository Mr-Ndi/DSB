// main.go
package main

import (
	"context"
	"dsb/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("DB_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("We are connected to mongoDB")

	UserCollection = client.Database("digitalsuggestionboxpro").Collection("coffee")
	r := mux.NewRouter()

	// Pass UserCollection to routes
	routes.UserRoutes(r, UserCollection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}

	log.Printf("Server will start at http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
