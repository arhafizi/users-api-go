package middlewares

import (
	"time"

	"example.com/api/pkg/logging"
	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger logging.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		logger.Info(logging.RequestResponse, logging.Api, "Request processed", map[logging.ExtraKey]any{
			logging.Method:     c.Request.Method,
			logging.Path:       c.Request.URL.Path,
			logging.ClientIp:   c.ClientIP(),
			logging.StatusCode: c.Writer.Status(),
			logging.Latency:    duration.String(),
		})
	}
}
