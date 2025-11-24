package main

import (
	"log"
	"monetix-be-api/configs"
	"monetix-be-api/internal/handlers"
	"monetix-be-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Set Gin mode
	gin.SetMode(config.Server.Mode)

	// Initialize database
	db, err := database.NewPostgresDB(&config.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize router with middleware
	router := gin.Default()
	router.Use(gin.Recovery())

	// Setup routes
	handlers.SetupUserRoutes(router, db)

	// Start server
	log.Printf("Server starting on %s in %s mode", config.Server.Port, config.Server.Mode)
	router.Run(config.Server.Port)
}
