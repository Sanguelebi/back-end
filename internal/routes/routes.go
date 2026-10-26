package routes

import (
	"course-selection-platform/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// The main SetupRouter function.
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "To access the api, add /api and the endpoint you are looking for.",
		})
	})

	{
		api := r.Group("/api")
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Pong!",
			})
		})
		api.GET("/subjects", handlers.GetSubjects(db))
	}

	return r;
}