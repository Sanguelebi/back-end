package handlers

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSubjectStudentOverviewHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		idVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized from access."})
			return
		}

		userId := idVal.(uint)

		status, data, err := services.GetSubjectStudentOverviewService(app, userId)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(status, gin.H{"message": "Subject Student overview succesfully fetched", "overview": data})
	}
}