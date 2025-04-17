package routes

import (
	"time"

	"example.com/api/internal/api/handlers"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
)

var store = persistence.NewInMemoryStore(time.Second)

func SetupUserRoutes(router *gin.RouterGroup, h *handlers.UserHandler) {
	users := router.Group("/users")
	{
		users.GET("", h.GetAll)
		users.GET("/cached", cache.CachePage(store, time.Minute, h.GetAll))
		users.GET("/:id", h.GetByID)
		users.POST("", h.Create)
		users.PUT("/:id", h.UpdateFull)
		users.PATCH("/:id", h.UpdatePartial)
		users.DELETE("/:id", h.DeleteUser)
	}
}
