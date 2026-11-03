package services

import (
	"course-selection-platform/internal/models"
	"errors"

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

func GetStudentPriorities(db *gorm.DB, studentID uint) ([]models.StudentPriorityList, error) {
	var priorities []models.StudentPriorityList

	err := db.Where("student_id = ?", studentID).
		Order("priority ASC").
		Find(&priorities).Error

	return priorities, err
}

func AddPriority(db *gorm.DB, studentID, groupID uint) error {
	// Verify that the SubjectGroup exists before inserting
	var group models.SubjectGroup
	if err := db.First(&group, groupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject group not found")
		}
		return err
	}

	// Prevent adding the exact same subject group twice
	var duplicate models.StudentPriorityList
	if err := db.Where("student_id = ? AND subject_group_id = ?", studentID, groupID).First(&duplicate).Error; err == nil {
		return errors.New("this subject group is already in your priority list")
	}

	// Check current priority count limit (Max 3)
	var count int64
	if err := db.Model(&models.StudentPriorityList{}).Where("student_id = ?", studentID).Count(&count).Error; err != nil {
		return err
	}

	if count >= 3 {
		return errors.New("maximum 3 priorities allowed")
	}

	// Automatically set priority (1, 2, or 3)
	entry := models.StudentPriorityList{
		StudentId:      studentID,
		SubjectGroupId: groupID,
		Priority:       int16(count + 1),
	}

	return db.Create(&entry).Error
}

func DeletePriority(db *gorm.DB, studentID, priorityID uint) error {
	res := db.Where("id = ? AND student_id = ?", priorityID, studentID).Delete(&models.StudentPriorityList{})
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New("priority item not found")
	}

	// Fetch remaining items and fix sequence numbers (1, 2...)
	var remaining []models.StudentPriorityList
	db.Where("student_id = ?", studentID).Order("priority asc").Find(&remaining)

	for i, item := range remaining {
		db.Model(&item).Update("priority", i+1)
	}

	return nil
}