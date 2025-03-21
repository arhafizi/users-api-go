package routes

import (
	"example.com/api/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, h *handlers.UserHandler) {
	users := router.Group("/users")
	{
		users.GET("", h.GetAll)
		users.GET("/:id", h.GetByID)
		users.POST("", h.Create)
		users.PUT("/:id", h.UpdateFull)
		users.PATCH("/:id", h.UpdatePartial)
		users.DELETE("/:id", h.DeleteUser)
	}
}
