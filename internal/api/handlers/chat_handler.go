package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"example.com/api/internal/api/responses"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services"
	"example.com/api/internal/services/chat"
	"example.com/api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	hub     *chat.Hub
	service services.IServiceManager
	logger  logging.ILogger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewChatHandler(h *chat.Hub, s services.IServiceManager, l logging.ILogger) *ChatHandler {
	return &ChatHandler{
		hub:     h,
		service: s,
		logger:  l,
	}
}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	userID, err := strconv.Atoi(c.GetString("user_id"))
	if err != nil {
		responses.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.service.User().GetByID(c.Request.Context(), int32(userID))
	if err != nil {
		responses.NotFound(c, "User not found")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error(logging.Internal, logging.Api, "Failed to upgrade connection", nil)
		return
	}

	client := chat.NewClient(h.hub, conn, int32(userID), user.Username)

	h.hub.Register(client)

	go client.SendMessages()
	go client.HandleMessages(h.service.Chat())
}

func (h *ChatHandler) GetMessageHistory(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cacheKey := fmt.Sprintf("messages:limit:%d:offset:%d", limit, offset)

	var messages []dbCtx.GetMessagesRow
	found, err := h.service.CacheStorage().Get(c.Request.Context(), cacheKey, &messages)
	if err != nil {
		slog.Error("cache", "redis", "Failed to get cache")
		responses.InternalServerError(c, "Failed to fetch message history")
		return
	}

	if found {
		responses.OK(c, "Message history retrieved from cache", messages)
		return
	}

	messages, err = h.service.Chat().GetMessages(c.Request.Context(), int32(limit), int32(offset))
	if err != nil {
		responses.InternalServerError(c, "Failed to fetch message history")
		return
	}

	if err := h.service.CacheStorage().Set(c.Request.Context(), cacheKey, messages, 5*time.Minute); err != nil {
		slog.Error("cache", "redis", "Failed to SET cache")
		responses.InternalServerError(c, "Failed to cache message history")
		return
	}
	responses.OK(c, "Message history retrieved successfully", messages)
}