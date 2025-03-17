package routes

import (
	"example.com/api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, handler *handlers.AuthHandler) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/register", handler.Register)
		auth.POST("/refresh", handler.Refresh)
	}
}
