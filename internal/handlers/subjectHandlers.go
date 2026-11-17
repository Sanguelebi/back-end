package handlers

import (
	"net/http"
	"strconv"

	"course-selection-platform/internal/models"
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

func GetStudentPrioritiesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized from access."})
			return
		}

		studentID := val.(uint)

		priorities, err := services.GetStudentPrioritiesService(db, studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Priorities fetched succesfully.", "priorities": priorities})
	}
}

func AddStudentPriorityHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		studentID := val.(uint)

		var req models.PriorityRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		if err := services.AddStudentPriorityService(db, studentID, req.SubjectGroupID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Priority added successfully"})
	}
}

func DeleteStudentPriorityHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.MustGet("userID").(uint)
		id, _ := strconv.Atoi(c.Param("id"))

		if err := services.DeleteStudentPriorityService(db, studentID, uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Priority deleted and re-indexed"})
	}
}

