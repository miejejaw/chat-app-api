package main

import (
	"chat-app-api/internal/database"
	"chat-app-api/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {

	// Check if the environment is development or production
	env := os.Getenv("GIN_ENV")
	if env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using Render's environment variables")
		}
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},                   // Allow specific origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // Allow specific HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow specific headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Allow credentials (cookies, authorization headers)
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	routes.SetupRoutes(api, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
