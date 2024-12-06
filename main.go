package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"dsb/config"
	_ "dsb/docs"
	"dsb/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RequestLogger logs incoming requests and their durations
func RequestLogger(c *gin.Context) {
	start := time.Now()
	log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
	c.Next() // Process the next handler
	duration := time.Since(start)
	log.Printf("Response sent: %s %s - %v", c.Request.Method, c.Request.URL.Path, duration)
}

func main() {
	// Load configuration from environment variables and handle errors
	if err := config.DatabaseConnection(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err) // Log error if connection fails
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://ds-box.onrender.com"}, // Replace with allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.UserRoutes(r) // Set up user routes
	routes.SuggestionRoutes(r)

	r.Use(RequestLogger) // Use the request logger middleware

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger documentation route

	port := os.Getenv("PORT") // Get PORT from environment variables
	if port == "" {
		port = "8080" // Default to port 8080 if not set
	}

	fmt.Printf("\n\n=======================================================>\n   Starting server on port %v\n=======================================================\n\n\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err) // Log fatal error if server fails to start
	}
}
