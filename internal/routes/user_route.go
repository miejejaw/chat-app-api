package routes

import (
	"chat-app-api/internal/handlers"
	"chat-app-api/internal/middleware"
	"chat-app-api/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, userService services.UserService) {
	userController := handlers.NewUserHandler(userService)

	userRoutes := router.Group("")
	{
		userRoutes.POST("/signup", userController.CreateUser)
		userRoutes.Use(middleware.AuthMiddleware())

		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/", userController.GetAllUsers)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
		userRoutes.GET("/search", userController.SearchUser)
	}
}
