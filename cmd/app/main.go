package main

import (
	"chat-app-api/internal/database"
	"chat-app-api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	router := gin.Default()
	api := router.Group("/api")
	routes.SetupRoutes(api, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
