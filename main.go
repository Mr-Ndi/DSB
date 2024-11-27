package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"dsb/config"
	"dsb/routes"

	_ "dsb/docs" // Import the generated swagger docs (if you're using swaggo)

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	// "github.com/swaggo/http-swagger" // Swagger UI
	// You may also need other imports based on your setup
)

// RequestLogger is a custom middleware for logging HTTP requests
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		// Record the start time
		start := time.Now()

		// Pass the request to the next handler
		next.ServeHTTP(w, r)

		// Log the response status and the duration of the request
		duration := time.Since(start)
		log.Printf("Response sent: %s %s - %v", r.Method, r.URL.Path, duration)
	})
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set up the database connection
	config.DatabaseConnection()

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Register routes
	routes.UserRoutes(r)

	// Serve Swagger UI at /swagger route
	r.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger-ui/"))))
	r.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger.yaml")
	})

	// Use the custom logger middleware
	r.Use(RequestLogger)

	// Get the port from the environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %v\n", port)

	// Start the HTTP server
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
