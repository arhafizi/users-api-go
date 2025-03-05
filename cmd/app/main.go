package main

import (
	"database/sql"
	"log"

	"example.com/api/internal/api/handlers"
	"example.com/api/internal/api/routes"
	"example.com/api/internal/repository"
	"example.com/api/internal/services"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db := initDb()

	// Initialize layers
	userRepo := repository.NewUserRepo(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()

	routes.SetupUserRoutes(router, userHandler)

	if err := router.Run(":5000"); err != nil {
		log.Fatal(err)
	}
}

func initDb() *sql.DB {
	conStr := "user=postgres password=4321 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)

	}
	if err := db.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
		defer db.Close()
	}
	return db
}
