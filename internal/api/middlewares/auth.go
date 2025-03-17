package middlewares

import (
	"errors"
	"strings"

	"example.com/api/internal/api/responses"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			responses.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			responses.Unauthorized(c, "Invalid token format")
			c.Abort()
			return
		}

		tokenString := bearerToken[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				println("here :::: ", err.Error())
				responses.Unauthorized(c, "Token expired")
			} else {
				responses.Unauthorized(c, "Invalid token")
			}
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			responses.Unauthorized(c, "Invalid token claims")
			c.Abort()
			return
		}

		userID, exists := claims["sub"].(string)
		if !exists {
			responses.Unauthorized(c, "Missing sub claim")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
