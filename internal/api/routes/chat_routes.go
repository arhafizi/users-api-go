package routes

import (
	"example.com/api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.RouterGroup, handler *handlers.ChatHandler) {
	chat := router.Group("/chat")
	{
		chat.GET("/ws", handler.HandleWebSocket)
		chat.GET("/messages", handler.GetMessageHistory)
	}
}
