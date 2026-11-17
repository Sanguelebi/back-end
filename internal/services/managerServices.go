package services

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func GetSubjectStudentOverviewService(app *config.AppData, managerId uint) (int, []models.SubjectStudentOverviewResponse, error) {
	var manager models.Users
	err := app.DB.Select("id", "faculty_id").First(&manager, managerId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusOK, make([]models.SubjectStudentOverviewResponse, 0), nil
		}
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch manager: %w", err)
	}

	if manager.FacultyID == nil {
		return http.StatusBadRequest, nil, fmt.Errorf("manager does not have a faculty assigned")
	}

	// Fetch all subjects belonging to the manager's faculty
	var subjects []models.Subject
	if err := app.DB.Where("faculty_id = ?", *manager.FacultyID).Find(&subjects).Error; err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch faculty subjects: %w", err)
	}

	if len(subjects) == 0 {
		return http.StatusOK, make([]models.SubjectStudentOverviewResponse, 0), nil
	}

	subjectIDs := make([]uint, len(subjects))
	for i, s := range subjects {
		subjectIDs[i] = s.ID
	}

	//  Fetch ALL groups for those subjects (including groups with 0 enrollments)
	var subjectGroups []models.SubjectGroup
	if err := app.DB.Where("subject_id IN ?", subjectIDs).Order("group_no asc").Find(&subjectGroups).Error; err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch subject groups: %w", err)
	}

	// Fetch active enrollments for those groups
	studentsByGroup := make(map[uint][]uint)
	if len(subjectGroups) > 0 {
		groupIDs := make([]uint, len(subjectGroups))
		for i, g := range subjectGroups {
			groupIDs[i] = g.ID
		}

		var enrollments []models.EnrollmentTable
		err := app.DB.
			Select("subject_group_id", "student_id").
			Where("subject_group_id IN ? AND status = ?", groupIDs, models.StatusActive).
			Find(&enrollments).Error
		if err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch enrollments: %w", err)
		}

		for _, e := range enrollments {
			studentsByGroup[e.SubjectGroupID] = append(studentsByGroup[e.SubjectGroupID], e.StudentID)
		}
	}

	// Group the SubjectGroupOverview items by SubjectID
	groupsBySubject := make(map[uint][]models.SubjectGroupOverview)
	for _, g := range subjectGroups {
		enrolledStudents := studentsByGroup[g.ID]
		if enrolledStudents == nil {
			enrolledStudents = make([]uint, 0) // Ensure JSON serializes as [] instead of null
		}

		groupsBySubject[g.SubjectId] = append(groupsBySubject[g.SubjectId], models.SubjectGroupOverview{
			GroupID:          g.ID,
			GroupNo:          g.GroupNo,
			StudentIDs:       enrolledStudents,
			AmountOfStudents: len(enrolledStudents),
			MaxAmount:        g.MaxStudent,
		})
	}

	// Assemble the final response
	response := make([]models.SubjectStudentOverviewResponse, 0, len(subjects))
	for _, subj := range subjects {
		groups := groupsBySubject[subj.ID]
		if groups == nil {
			groups = make([]models.SubjectGroupOverview, 0)
		}

		response = append(response, models.SubjectStudentOverviewResponse{
			Subject:       subj,
			SubjectGroups: groups,
		})
	}

	return http.StatusOK, response, nil
}