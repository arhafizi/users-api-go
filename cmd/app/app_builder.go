package main

import (
	"database/sql"

	"example.com/api/config"
	dbCfg "example.com/api/db"
	"example.com/api/internal/api/middlewares"
	"example.com/api/pkg/logging"
	"github.com/gin-gonic/gin"
)

type AppBuilder struct {
	Config *config.Config
	Logger logging.ILogger
	DB     *sql.DB
	Router *gin.Engine
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func (b *AppBuilder) WithConfig() *AppBuilder {
	b.Config = config.GetConfig()
	return b
}

func (b *AppBuilder) WithLogger() *AppBuilder {
	b.Logger = logging.NewLogger(b.Config)
	return b
}

func (b *AppBuilder) WithDB() *AppBuilder {
	b.DB = dbCfg.InitDb(b.Config, b.Logger)
	return b
}

func (b *AppBuilder) WithRouter() *AppBuilder {
	b.Router = gin.New()
	b.Router.Use(gin.Recovery())
	b.Router.Use(middlewares.LoggingMiddleware(b.Logger))
	return b
}

func (b *AppBuilder) Build() *gin.Engine {
	return b.Router
}
