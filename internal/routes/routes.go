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

	// Set up services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	// Authentication routes
	authRoutes := router.Group("/auth")
	SetupAuthRoutes(authRoutes, authService)

	// User routes
	userRoutes := router.Group("/users")
	SetupUserRoutes(userRoutes, userService)
}
