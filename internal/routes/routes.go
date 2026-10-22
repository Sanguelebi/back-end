package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// The main SetupRouter function.
func SetupRouter() *gin.Engine {
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
	}

	return r;
}