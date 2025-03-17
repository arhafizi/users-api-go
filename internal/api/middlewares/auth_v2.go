package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"example.com/api/internal/api/responses"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthV2Middleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookie, err := c.Cookie("token"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			responses.Unauthorized(c, "Authentication required: no token provided")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
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

		var userID string

		if sub, exists := claims["sub"]; exists {
			if uid, ok := sub.(string); ok {
				userID = uid
			} else {
				responses.Unauthorized(c, "Invalid token: sub claim format")
				c.Abort()
				return
			}
		} else if claimUserID, exists := claims["user_id"]; exists {
			if uid, ok := claimUserID.(string); ok {
				userID = uid
			} else {
				responses.Unauthorized(c, "Invalid token: user_id claim format")
				c.Abort()
				return
			}
		} else {
			responses.Unauthorized(c, "Missing user identifier in token claim")
			c.Abort()
			return
		}

		// db lookup
		// user, err := store.GetUserByID(context.Background(), userID)
		// if err != nil {
		// 	responses.Forbidden(c, "User not found or unauthorized")
		// 	c.Abort()
		// 	return
		// }

		c.Set("user_id", userID)
		c.Next()
	}
}
