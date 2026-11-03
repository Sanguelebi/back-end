package handlers

import (
	"net/http"

	"course-selection-platform/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSubjects(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		faculty := c.Query("faculty")
		name := c.Query("name")

		subjects, err := services.GetAllSubjects(db, faculty, name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, subjects)
	}
}