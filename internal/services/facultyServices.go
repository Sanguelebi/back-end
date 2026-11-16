package services

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
)

func GetFacultiesService(app *config.AppData) ([]models.Faculty, error) {
	var faculties []models.Faculty
	result := app.DB.Find(&faculties)
	return faculties, result.Error
}
