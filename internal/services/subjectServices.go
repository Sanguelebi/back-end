package services

import (
	"course-selection-platform/internal/models"

	"gorm.io/gorm"
)

func GetAllSubjects(db *gorm.DB) ([]models.Subject, error) {
	var subjects []models.Subject
	err := db.Find(&subjects).Error
	return subjects, err
}