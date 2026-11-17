package handlers

import (
	"net/http"

	"course-selection-platform/config"
	"course-selection-platform/internal/services"

	"github.com/gin-gonic/gin"
)

func GetFacultiesHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		faculties, err := services.GetFacultiesService(app)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Got faculties succesfully!", "faculties": faculties})
	}
}
