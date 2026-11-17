package handlers

import (
	"net/http"
	"strconv"

	"course-selection-platform/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllSubjectsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		faculty := c.Query("faculty")
		name := c.Query("name")

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		subjects, hasMore, err := services.GetAllSubjectsService(db, faculty, name, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Got subjects succesfully!", "subjects": subjects, "has_more": hasMore})
	}
}