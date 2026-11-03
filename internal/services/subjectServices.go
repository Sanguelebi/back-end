package services

import (
	"course-selection-platform/internal/models"

	"gorm.io/gorm"
)

func GetAllSubjects(db *gorm.DB, facultyID, name string) ([]models.Subject, error) {
	var subjects []models.Subject
	query := db.Model(&models.Subject{})

	if facultyID != "" {
		query = query.Where("faculty_id = ?", facultyID)
	}

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	err := query.Find(&subjects).Error
	return subjects, err
}