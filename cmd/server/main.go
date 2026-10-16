package main

import (
	"course-selection-platform/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Course Selection Platform API",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()

	log.Println("Database connected successfully!")

}
