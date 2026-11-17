package handlers

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CoordinatorConfirmSubjectHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		idQuery := c.Query("id")
		if idQuery == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No priority ID Specified."})
			return
		}

		convertPriorityId, err := strconv.Atoi(idQuery)
		if err != nil || convertPriorityId < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		priorityId := uint(convertPriorityId)
		status, err := services.CoordinatorConfirmSubjectService(app, priorityId)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(status, gin.H{"message": "Coordinator changed confirm succesfully!"})
	}
}

func GetCoordinatorStudentPrioritiesHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		idVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized from access."})
			return
		}

		userId := idVal.(uint)

		status, data, err := services.GetCoordinatorStudentPrioritiesService(app, userId)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
		}
		c.JSON(status, gin.H{"message": "Student priorities fetched succesfully!", "data": data})
	}
}