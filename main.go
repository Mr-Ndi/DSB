package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"dsb/config"
	_ "dsb/docs"
	"dsb/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RequestLogger logs incoming requests and their durations
func RequestLogger(c *gin.Context) {
	start := time.Now()
	log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
	c.Next()
	duration := time.Since(start)
	log.Printf("Response sent: %s %s - %v", c.Request.Method, c.Request.URL.Path, duration)
}

// loadEnv loads environment variables from a .env file if not in production
func loadEnv() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}
}

func main() {
	loadEnv() // Load environment variables

	config.DatabaseConnection() // Initialize database connection

	r := gin.Default()

	routes.UserRoutes(r) // Set up routes

	r.Use(RequestLogger) // Use request logger middleware

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger documentation route

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}

	fmt.Printf("Starting server on port %v\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
