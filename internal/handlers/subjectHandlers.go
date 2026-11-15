package handlers

import (
	"net/http"
	"strconv"

	"course-selection-platform/internal/models"
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

		c.JSON(http.StatusOK, gin.H{"message": "Got subjects succesfully!", "subjects": subjects})
	}
}

func GetPriorities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized from access."})
			return
		}

		studentID := val.(uint)

		priorities, err := services.GetStudentPriorities(db, studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Priorities fetched succesfully.", "priorities": priorities})
	}
}

func AddPriority(db *gorm.DB) gin.HandlerFunc {
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

		if err := services.AddPriority(db, studentID, req.SubjectGroupID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Priority added successfully"})
	}
}

func DeletePriority(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.MustGet("userID").(uint)
		id, _ := strconv.Atoi(c.Param("id"))

		if err := services.DeletePriority(db, studentID, uint(id)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Priority deleted and re-indexed"})
	}
}