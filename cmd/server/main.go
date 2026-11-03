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
	
	dbCfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("DB Config Initialization failed: %v", err)
	}
	
	db, err := database.ConnectDB(dbCfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Load the JWT Secret from .env
	appData, err := config.LoadAppConfig(db)
	if err != nil {
		log.Fatalf("App Config initialization has failed: %v", err)
	}
	_ = appData

	// This SetupRouter() function will contain all of our API endpoints.
	router := routes.SetupRouter(appData)

	// Get PORT from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // If can't access, set a default.
	}

	// Use that PORT here.
	router.Run(":" + port)
}
