package main

import (
	"course-selection-platform/config"
	"course-selection-platform/database"
	"course-selection-platform/internal/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env values first
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return
	}
	
	// Load the DB Config first and then connect to the DB with it.
	cfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Config initialization has failed: %v", err)
	}
	
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	_ = db // This line here is to avoid getting yelled at by the compiler.

	// This SetupRouter() function will contain all of our API endpoints.
	router := routes.SetupRouter()
	
	// Get PORT from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // If can't access, set a default.
	}

	// Use that PORT here.
	router.Run(":" + port)
}
