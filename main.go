package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"dsb/config"
	"dsb/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	config.DatabaseConnection()

	r := mux.NewRouter()
	routes.UserRoutes(r)

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file:%v", err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %v\n", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
