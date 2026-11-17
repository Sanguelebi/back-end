package services

import (
	"course-selection-platform/internal/models"

	"gorm.io/gorm"
)

func GetAllSubjectsService(db *gorm.DB, facultyID, name string, page, pageSize int) ([]models.Subject, bool, error) {
	// === Defaults ===
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	var subjects []models.Subject
	query := db.Model(&models.Subject{})

	if facultyID != "" {
		query = query.Where("faculty_id = ?", facultyID)
	}
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	offset := (page - 1) * pageSize

	// Fetch pageSize + 1 to check if another page exists
	err := query.
		Order("id ASC").
		Offset(offset).
		Limit(pageSize + 1).
		Find(&subjects).Error
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if len(subjects) > pageSize {
		hasMore = true
		subjects = subjects[:pageSize] // Drop the extra probe item
	}

	return subjects, hasMore, nil
}