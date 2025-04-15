package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(5, 20)

// var limiter = rate.NewLimiter(1, 5)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Too many requests"})
			return
		}
		c.Next()
	}
}
