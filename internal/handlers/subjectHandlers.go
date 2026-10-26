package handlers

import (
	"net/http"

	"course-selection-platform/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSubjects(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		subjects, err := services.GetAllSubjects(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, subjects)
	}
}