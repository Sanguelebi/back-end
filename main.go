package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
