package middlewares

import (
	"time"

	"fmt"

	"example.com/api/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(metrics.HttpDuration)
	prometheus.MustRegister(metrics.DbCall)
	prometheus.MustRegister(metrics.TotalReq)
	prometheus.MustRegister(metrics.NodeUsage)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		metrics.TotalReq.Inc()
		c.Next()

		duration := time.Since(start).Seconds()

		// Record the duration
		metrics.HttpDuration.WithLabelValues(
			c.Request.URL.Path,
			c.Request.Method,
			fmt.Sprint(c.Writer.Status()),
		).Observe(duration)
	}
}
