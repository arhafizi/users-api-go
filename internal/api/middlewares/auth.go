package middlewares

import (
	"strings"

	"example.com/api/internal/api/responses"
	"example.com/api/internal/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			responses.Unauthorized(c, "Authentication required: no token provided")
			c.Abort()
			return
		}

		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			responses.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			responses.Unauthorized(c, "Invalid token: missing or invalid sub claim")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if authHeader := c.GetHeader("Authorization"); authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
		return cookie
	}

	return ""
}
