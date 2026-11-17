package handlers

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
	"course-selection-platform/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStudentChoosableSubjectsHandler(app *config.AppData) gin.HandlerFunc {
	return func(c *gin.Context) {
		idVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized from access."})
			return
		}

		facultyQuery := c.Query("subjectType")
		if facultyQuery == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Specify subject type."})
			return
		}
		faculty, err := models.ParseSubjectType(facultyQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect subject type."})
			return
		}

		studentID := idVal.(uint)
		status, data, err := services.GetStudentChoosableSubjectsService(app, studentID, faculty)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
		c.JSON(status, gin.H{"message": "Choosable Subjects Fetched Succesfully!", "subjects": data})
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