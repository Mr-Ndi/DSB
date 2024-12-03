package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"dsb/config"
	_ "dsb/docs" // Import your docs package for swagger
	"dsb/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RequestLogger(c *gin.Context) {
	start := time.Now()
	log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
	c.Next()
	duration := time.Since(start)
	log.Printf("Response sent: %s %s - %v", c.Request.Method, c.Request.URL.Path, duration)
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set up database connection
	config.DatabaseConnection()

	// Initialize Gin router
	r := gin.Default()

	// Register routes
	routes.UserRoutes(r)

	// Apply request logger middleware
	r.Use(RequestLogger)

	// Serve Swagger UI (ensure this is only registered once)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("Starting server on port %v\n", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
