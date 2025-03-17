package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/api/internal/api/responses"
	"example.com/api/internal/api/validation"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services"
	"example.com/api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service services.IServiceManager
	logger  logging.ILogger
}

func NewUserHandler(s services.IServiceManager, logger logging.ILogger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  logger,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dbCtx.CreateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(logging.Validation, logging.Api, "Invalid request body", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
			logging.Path:         c.Request.URL.Path,
			logging.Method:       c.Request.Method,
		})

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			validationErrors := validation.GetValidationErrors(err)
			responses.BadRequest(c, "Invalid request body", validationErrors)
			return
		}

		responses.BadRequest(c, "Invalid request body", err)
		return
	}

	user, err := h.service.User().Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Error(logging.Internal, logging.FailedToCreateUser, "Failed to create user", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
			logging.Path:         c.Request.URL.Path,
			logging.Method:       c.Request.Method,
		})
		responses.InternalServerError(c, "Failed to create user")
		return
	}

	h.logger.Info(logging.Internal, logging.Api, "User created successfully", map[logging.ExtraKey]any{
		logging.Path:   c.Request.URL.Path,
		logging.Method: c.Request.Method,
	})

	responses.Created(c, "User created successfully", user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, "Invalid user ID, must be an integer", nil)
		return
	}

	user, err := h.service.User().GetByID(c.Request.Context(), int32(id))
	if err != nil {
		responses.NotFound(c, "User not found")
		return
	}

	responses.OK(c, "User retrieved successfully", user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		responses.BadRequest(c, "Invalid limit parameter", nil)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		responses.BadRequest(c, "Invalid offset parameter", nil)
		return
	}

	users, err := h.service.User().GetAll(c.Request.Context(), dbCtx.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		responses.InternalServerError(c, "Failed to retrieve users")
		return
	}

	responses.OK(c, "Users retrieved successfully", users)
}

func (h *UserHandler) UpdateFull(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, "Invalid user ID, must be an integer", nil)
		return
	}

	var req dbCtx.UpdateUserFullParams
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Invalid request body", err)
		return
	}

	req.ID = int32(id)
	user, err := h.service.User().UpdateFull(c.Request.Context(), req)
	if err != nil {
		if user == nil {
			responses.NotFound(c, "User not found")
			return
		}
		responses.InternalServerError(c, "Failed to update user")
		return
	}

	responses.OK(c, "User updated successfully", user)
}

func (h *UserHandler) UpdatePartial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, "Invalid user ID, must be an integer", nil)
		return
	}

	var req dbCtx.UpdateUserPartialParams
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Invalid request body", err)
		return
	}

	req.ID = int32(id)
	user, err := h.service.User().UpdatePartial(c.Request.Context(), req)
	if err != nil {
		if user == nil {
			responses.NotFound(c, "User not found")
			return
		}
		responses.InternalServerError(c, "Failed to partially update user")
		return
	}

	responses.OK(c, "User partially updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, "Invalid user ID, must be an integer", nil)
		return
	}

	err = h.service.User().SoftDelete(c.Request.Context(), int32(id))
	if err != nil {
		responses.InternalServerError(c, "Failed to delete user")
		return
	}

	c.Status(http.StatusNoContent)
}
