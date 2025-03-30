package handlers

import (
	"errors"
	"fmt"

	"example.com/api/internal/api/responses"
	"example.com/api/internal/api/validation"
	dto "example.com/api/internal/contracts"
	"example.com/api/internal/services"
	"example.com/api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService services.IAuthService
	logger      logging.ILogger
}

func NewAuthHandler(s services.IAuthService, l logging.ILogger) *AuthHandler {
	return &AuthHandler{
		authService: s,
		logger:      l,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			validationErrors := validation.GetValidationErrors(err)
			responses.BadRequest(c, "Invalid request body", validationErrors)
			return
		}
		responses.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	user, err := h.authService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		// responses.Unauthorized(c, "Invalid credentials")
		responses.Unauthorized(c, err.Error())
		return
	}

	accessToken, err := h.authService.GenerateAccessToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		h.logger.Error(logging.General, logging.SubCategory(logging.Internal), "Failed to generate access token", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
			logging.Path:         c.Request.URL.Path,
			logging.Method:       c.Request.Method,
		})
		responses.InternalServerError(c, "Failed to generate access token")
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		h.logger.Error(logging.General, logging.SubCategory(logging.Internal), "Failed to generate refresh token", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
			logging.Path:         c.Request.URL.Path,
			logging.Method:       c.Request.Method,
		})
		responses.InternalServerError(c, "Failed to generate refresh token")
		return
	}

	responses.OK(c, "User logged in successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Invalid request body", err)
		return
	}

	user, accessToken, refreshToken, err := h.authService.RegisterUser(c.Request.Context(), req)
	if err != nil {
		h.logger.Error(logging.Internal, logging.FailedToCreateUser, "Failed to register user", map[logging.ExtraKey]any{
			logging.ErrorMessage: err.Error(),
			logging.Path:         c.Request.URL.Path,
			logging.Method:       c.Request.Method,
		})
		responses.BadRequest(c, err.Error(), nil)
		return
	}

	responses.Created(c, "User registered successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "Refresh token is required", err)
		return
	}

	accessToken, refreshToken, err := h.authService.RotateTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		responses.Unauthorized(c, err.Error())
		return
	}

	responses.OK(c, "Tokens refreshed successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
