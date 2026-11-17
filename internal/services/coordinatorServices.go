package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"course-selection-platform/config"
	"course-selection-platform/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CoordinatorConfirmSubjectService(app *config.AppData, priorityId uint) (int, error) {
	statusCode := http.StatusInternalServerError

	err := app.DB.Transaction(func(tx *gorm.DB) error {
		var priority models.StudentPriorityList
		// Lock it down so that the other coordinator doesn't mess anything up.
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("SubjectGroup").
			First(&priority, priorityId).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
				return fmt.Errorf("priority record with id %d not found", priorityId)
			}
			return fmt.Errorf("failed to fetch priority record: %w", err)
		}

		if !priority.Confirm {
            var activeCount int64
            if err := tx.Model(&models.EnrollmentTable{}).
                Where("subject_group_id = ? AND status = ?", priority.SubjectGroupId, models.StatusActive).
                Count(&activeCount).Error; err != nil {
                return fmt.Errorf("failed to check group capacity: %w", err)
            }

            if activeCount >= int64(priority.SubjectGroup.MaxStudent) {
                statusCode = http.StatusConflict
                return fmt.Errorf("subject group %d is full (%d/%d)", priority.SubjectGroupId, activeCount, priority.SubjectGroup.MaxStudent)
            }

            currentSemester, err := calculateCurrentSemester(tx, priority.StudentId)
            if err != nil {
                return fmt.Errorf("failed to calculate student semester: %w", err)
            }

            // Find and unconfirm any existing confirmed priority in this exact slot
            var existingConfirmed models.StudentPriorityList
            err = tx.Where("student_id = ? AND slot_id = ? AND confirm = ? AND id != ?", 
                priority.StudentId, priority.SlotID, true, priority.ID).
                First(&existingConfirmed).Error

            if err == nil {
                // If one exists, delete its enrollment
                if err := tx.Where("student_id = ? AND subject_group_id = ? AND status = ?",
                    existingConfirmed.StudentId, existingConfirmed.SubjectGroupId, models.StatusActive,
                ).Delete(&models.EnrollmentTable{}).Error; err != nil {
                    return fmt.Errorf("failed to remove existing enrollment: %w", err)
                }
                // Mark it as unconfirmed
                if err := tx.Model(&existingConfirmed).Update("confirm", false).Error; err != nil {
                    return fmt.Errorf("failed to unconfirm previous priority: %w", err)
                }
            } else if !errors.Is(err, gorm.ErrRecordNotFound) {
                return fmt.Errorf("database error checking for existing confirmations: %w", err)
            }

            enrollment := models.EnrollmentTable{
                Status:         models.StatusActive,
                SubjectGroupID: priority.SubjectGroupId,
                Year:           int16(time.Now().Year()),
                Semester:       currentSemester,
                StudentID:      priority.StudentId,
            }

            if err := tx.FirstOrCreate(&enrollment, models.EnrollmentTable{
                StudentID:      priority.StudentId,
                SubjectGroupID: priority.SubjectGroupId,
                Status:         models.StatusActive,
            }).Error; err != nil {
                return fmt.Errorf("failed to create enrollment: %w", err)
            }

            if err := tx.Model(&priority).Update("confirm", true).Error; err != nil {
                return fmt.Errorf("failed to confirm priority record: %w", err)
            }
            priority.Confirm = true
        } else {
			if err := tx.Where(
				"student_id = ? AND subject_group_id = ? AND status = ?",
				priority.StudentId,
				priority.SubjectGroupId,
				models.StatusActive,
			).Delete(&models.EnrollmentTable{}).Error; err != nil {
				return fmt.Errorf("failed to remove enrollment: %w", err)
			}

			if err := tx.Model(&priority).Update("confirm", false).Error; err != nil {
				return fmt.Errorf("failed to unconfirm priority record: %w", err)
			}
			priority.Confirm = false
		}

		return nil
	})

	if err != nil {
		return statusCode, err
	}

	return http.StatusOK, nil
}

func GetCoordinatorStudentPrioritiesService(app *config.AppData, coordinatorID uint) (int, []models.StudentPriorityList, error) {
    var coordinator models.Users
    if err := app.DB.First(&coordinator, coordinatorID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return http.StatusNotFound, nil, fmt.Errorf("coordinator with id %d not found", coordinatorID)
        }
        return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch coordinator: %w", err)
    }

    if coordinator.FacultyID == nil {
        return http.StatusBadRequest, nil, errors.New("coordinator is not assigned to any faculty")
    }

    priorities := make([]models.StudentPriorityList, 0)
    err := app.DB.
        Preload("SubjectGroup"). // ADDED: Preload so the coordinator can see the subject details
        Joins("JOIN users ON users.id = student_priority_lists.student_id").
        Where("users.faculty_id = ? AND users.role = ?", *coordinator.FacultyID, models.RoleStudent).
        // ADDED: slot_id ASC sorting so alternatives are grouped together
        Order("student_priority_lists.student_id ASC, student_priority_lists.slot_id ASC, student_priority_lists.priority ASC").
        Find(&priorities).Error
    
    if err != nil {
        return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch student priorities: %w", err)
    }

    return http.StatusOK, priorities, nil
}