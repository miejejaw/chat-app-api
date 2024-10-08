package routes

import (
	"chat-app-api/internal/handlers"
	"chat-app-api/internal/middleware"
	"chat-app-api/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupMessageRoutes(router *gin.RouterGroup, messageService services.MessageService) {
	messageHandler := handlers.NewMessageHandler(messageService)

	messageRoutes := router.Group("/")
	{
		messageRoutes.GET("/ws", messageHandler.HandleConnections)
		messageRoutes.GET("/friends", middleware.AuthMiddleware(), messageHandler.GetFriendsWithLastMessage)
		messageRoutes.GET("/friend/chats", middleware.AuthMiddleware(), messageHandler.GetMessagesBySenderIdAndReceiverId)
	}
}
