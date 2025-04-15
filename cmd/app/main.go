package main

import (
	"log"
	"runtime"
	"time"

	"example.com/api/config"
	dbConf "example.com/api/db"
	"example.com/api/internal/api/handlers"
	"example.com/api/internal/api/middlewares"
	"example.com/api/internal/api/routes"
	"example.com/api/internal/repository"
	"example.com/api/internal/services"
	"example.com/api/internal/services/chat"
	"example.com/api/pkg/logging"
	"example.com/api/pkg/metrics"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	conf := config.GetConfig()
	logger := logging.NewLogger(conf)
	db := dbConf.InitDb(conf, logger)

	repoManager := repository.NewRepositoryManager(db)
	serviceManager := services.NewServiceManager(repoManager, logger, *conf)

	userHandler := handlers.NewUserHandler(serviceManager, logger)
	authHandler := handlers.NewAuthHandler(serviceManager.Auth(), logger)

	app := gin.New()
	app.Use(
		gin.Recovery(),
		middlewares.LoggingMiddleware(logger),
		middlewares.PrometheusMiddleware(),
		middlewares.CORS(),
		middlewares.RateLimiter(),
		middlewares.Secure,
	)

	app.SetTrustedProxies([]string{"127.0.0.1"})

	routes.SetupMetricsRoutes(app)
	routes.SetupAuthRoutes(app, authHandler)
	protected := app.Group("/api")
	protected.Use(middlewares.AuthMiddleware(serviceManager.Auth()))
	routes.SetupUserRoutes(protected, userHandler)

	hub := chat.NewHub()
	go hub.Run()

	chatHandler := handlers.NewChatHandler(hub, serviceManager, logger)
	routes.SetupChatRoutes(protected, chatHandler)

	monitorSystemMetrics()

	if err := app.Run(":5000"); err != nil {
		log.Fatal(err)
	}
}

func monitorSystemMetrics() {
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// Memory usage in MB
			metrics.NodeUsage.WithLabelValues("memory", "heap").Set(float64(m.HeapAlloc) / 1024 / 1024)
			metrics.NodeUsage.WithLabelValues("memory", "stack").Set(float64(m.StackInuse) / 1024 / 1024)
			metrics.NodeUsage.WithLabelValues("memory", "system").Set(float64(m.Sys) / 1024 / 1024)

			// Goroutines
			metrics.NodeUsage.WithLabelValues("goroutines", "count").Set(float64(runtime.NumGoroutine()))

			time.Sleep(10 * time.Second)
		}
	}()
}
