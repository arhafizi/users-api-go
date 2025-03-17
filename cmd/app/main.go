package main

import (
	"log"

	"example.com/api/config"
	dbConf "example.com/api/db"
	"example.com/api/internal/api/handlers"
	"example.com/api/internal/api/middlewares"
	"example.com/api/internal/api/routes"
	"example.com/api/internal/repository"
	"example.com/api/internal/services"
	"example.com/api/pkg/logging"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	conf := config.GetConfig()
	logger := logging.NewLogger(conf)
	db := dbConf.InitDb(conf, logger)

	repoManager := repository.NewRepositoryManager(db)
	serviceManager := services.NewServiceManager(repoManager, logger,conf.JWT)

	userHandler := handlers.NewUserHandler(serviceManager, logger)
	authHandler := handlers.NewAuthHandler(serviceManager.Auth(), logger)

	app := gin.New()
	app.Use(gin.Recovery(), middlewares.LoggingMiddleware(logger))
	app.SetTrustedProxies([]string{"127.0.0.1"})

	routes.SetupAuthRoutes(app, authHandler)
	protected := app.Group("/api")
	protected.Use(middlewares.AuthV2Middleware(conf.JWT.Secret))
	routes.SetupUserRoutes(protected, userHandler)

	if err := app.Run(":5000"); err != nil {
		log.Fatal(err)
	}
}
