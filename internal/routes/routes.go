package routes

import (
	"chat-app-api/internal/repositories"
	"chat-app-api/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Set up repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Set up services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	messageService := services.NewMessageService(messageRepo)

	// routes
	authRoutes := router.Group("/auth")
	userRoutes := router.Group("/users")
	messageRoutes := router.Group("/messages")

	// Setup routes
	SetupAuthRoutes(authRoutes, authService)
	SetupUserRoutes(userRoutes, userService)
	SetupMessageRoutes(messageRoutes, messageService)
}
