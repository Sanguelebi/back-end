package main

import (
	"course-selection-platform/config"
	"course-selection-platform/database"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Course Selection Platform API",
		})
	})
	
	// Get PORT from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // If can't access, set a default.
	}

	// Use that PORT here.
	router.Run(":" + port)
}
