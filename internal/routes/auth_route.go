package routes

import (
	"chat-app-api/internal/handlers"
	"chat-app-api/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.RouterGroup, authService services.AuthService) {
	authHandler := handlers.NewAuthHandler(authService)

	authRouter := router.Group("")
	{
		authRouter.POST("/login", authHandler.Login)
		authRouter.POST("/renew", authHandler.RenewAccessToken)
	}
}
