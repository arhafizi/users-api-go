package handlers

import (
	"net/http"
	"strconv"

	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	srv services.IUserService
}

func NewUserHandler(s services.IUserService) *UserHandler {
	return &UserHandler{
		srv: s,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dbCtx.CreateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	user, err := h.srv.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"status":  "success",
		"data":    user,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID, must be an integer",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	user, err := h.srv.GetByID(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User retrieved successfully",
		"status":  "success",
		"data":    user,
	})
}

func (h *UserHandler) GetAll(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid limit parameter",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid offset parameter",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	users, err := h.srv.GetAll(c.Request.Context(), dbCtx.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve users",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"status":  "success",
		"data":    users,
	})
}

func (h *UserHandler) UpdateFull(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID, must be an integer",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	var req dbCtx.UpdateUserFullParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	req.ID = int32(id)
	user, err := h.srv.UpdateFull(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"status":  "success",
		"data":    user,
	})
}

func (h *UserHandler) UpdatePartial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID, must be an integer",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	var req dbCtx.UpdateUserPartialParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	req.ID = int32(id)
	user, err := h.srv.UpdatePartial(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to partially update user",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User partially updated successfully",
		"status":  "success",
		"data":    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID, must be an integer",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	err = h.srv.SoftDelete(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete user",
			"status":  "fail",
			"data":    nil,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
