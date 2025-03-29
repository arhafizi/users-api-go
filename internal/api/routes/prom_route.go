package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupMetricsRoutes(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
