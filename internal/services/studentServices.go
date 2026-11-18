package services

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func GetStudentPrioritiesService(db *gorm.DB, studentID uint) ([]models.StudentPriorityList, error) {
    var priorities []models.StudentPriorityList

    err := db.Preload("SubjectGroup").
        Where("student_id = ?", studentID).
        Order("slot_id ASC, priority ASC").
        Find(&priorities).Error

    return priorities, err
}

func AddStudentPriorityService(db *gorm.DB, studentID, groupID uint, slotID int, subjectType models.SubjectType) error {
	var group models.SubjectGroup
	if err := db.First(&group, groupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject group not found")
		}
		return err
	}

	var duplicate models.StudentPriorityList
	if err := db.Where("student_id = ? AND subject_group_id = ?", studentID, groupID).First(&duplicate).Error; err == nil {
		return errors.New("this subject group is already in your priority list")
	}

	var count int64
	if err := db.Model(&models.StudentPriorityList{}).
		Where("student_id = ? AND slot_id = ?", studentID, slotID).
		Count(&count).Error; err != nil {
		return err
	}

	if count >= 3 {
		return errors.New("maximum 3 priorities allowed per slot")
	}

	entry := models.StudentPriorityList{
		StudentId:      studentID,
		SubjectGroupId: groupID,
		SlotID:         uint(slotID),
		SubjectType:    subjectType,
		Priority:       int16(count + 1),
	}

	return db.Create(&entry).Error
}


func DeleteStudentPriorityService(db *gorm.DB, studentID, priorityID uint) error {
	var target models.StudentPriorityList

	// Find the item first to extract its SlotID before deleting
	if err := db.Where("id = ? AND student_id = ?", priorityID, studentID).First(&target).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("priority item not found")
		}
		return err
	}

	slotID := target.SlotID

	// Delete the item
	if err := db.Delete(&target).Error; err != nil {
		return err
	}

	// Fetch remaining items ONLY in the affected slot and fix sequence numbers (1, 2...)
	var remaining []models.StudentPriorityList
	db.Where("student_id = ? AND slot_id = ?", studentID, slotID).Order("priority asc").Find(&remaining)

	// Reorder smoothly so there are no gaps
	for i, item := range remaining {
		db.Model(&item).Update("priority", int16(i+1))
	}

	return nil
}

func GetStudentSubjectAllowanceService(app *config.AppData, studentID uint) (int, models.StudentSubjectAllowance, error) {
	currentSemester, err := calculateCurrentSemester(app.DB, studentID)
	if err != nil {
		return http.StatusInternalServerError, models.StudentSubjectAllowance{}, fmt.Errorf("failed to calculate student semester: %w", err)
	}

	var allowance models.StudentSubjectAllowance
	if err := app.DB.
		Where("student_id = ? AND semester = ?", studentID, currentSemester).
		First(&allowance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, models.StudentSubjectAllowance{}, fmt.Errorf("allowance not found for student %d in semester %d", studentID, currentSemester)
		}
		return http.StatusInternalServerError, models.StudentSubjectAllowance{}, fmt.Errorf("failed to fetch student subject allowance: %w", err)
	}

	return http.StatusOK, allowance, nil
}

func GetStudentChoosableSubjectsService(app *config.AppData, studentID uint, subjectType models.SubjectType) (int, []models.ChoosableSubjectResponse, error) {
	// Dynamically calculate the student's current semester based on enrollment history
	currentSemester, err := calculateCurrentSemester(app.DB, studentID)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to calculate current semester: %w", err)
	}

	// Fetch the student
	var student models.Users
	if err := app.DB.Select("faculty_id").First(&student, studentID).Error; err != nil {
		return http.StatusInternalServerError, nil, err
	}

	// Safety check just in case a student somehow has no faculty assigned
	if student.FacultyID == nil {
		return http.StatusBadRequest, nil, fmt.Errorf("student is not assigned to a faculty")
	}

	// Fetch initially ChoosableSubjects within the calculated semester bounds
	var choosableSubjects []models.ChoosableSubjects
	err = app.DB.
		Where("subject_type = ?", subjectType).
		Where("from_semester <= ?", currentSemester).
		Where("to_semester >= ?", currentSemester).
		Where("faculty_id = ?", *student.FacultyID).
		Preload("Subject").
		Find(&choosableSubjects).Error

	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch choosable subjects: %w", err)
	}

	if len(choosableSubjects) == 0 {
		return http.StatusOK, []models.ChoosableSubjectResponse{}, nil
	}

	// Get passed subjects for prerequisite checking and filtering
	passedMap, err := getPassedSubjects(app.DB, studentID)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to fetch passed subjects: %w", err)
	}

	// Fetch prerequisites and evaluate valid subjects
	validChoosableSubjects, validSubjectIDs, err := evaluatePrerequisites(app.DB, choosableSubjects, passedMap)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to evaluate prerequisites: %w", err)
	}

	if len(validSubjectIDs) == 0 {
		return http.StatusOK, []models.ChoosableSubjectResponse{}, nil
	}

	// Fetch SubjectGroups and map to final response DTO
	response, err := buildSubjectResponseDTO(app.DB, validChoosableSubjects, validSubjectIDs)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, response, nil
}

// --- HELPER FUNCTIONS ---

func calculateCurrentSemester(db *gorm.DB, studentID uint) (models.Semester, error) {
    // Find the highest semester number the student has ever enrolled in
    var maxSemester models.Semester
    err := db.Model(&models.EnrollmentTable{}).
        Where("student_id = ?", studentID).
        Select("COALESCE(MAX(semester), 0)").
        Scan(&maxSemester).Error

    if err != nil {
        return 0, err
    }

    // Freshman with 0 enrollment records defaults to Semester 1
    if maxSemester == 0 {
        return 1, nil
    }

    // Check if the student still has ANY active (in-progress) subjects in that semester
    var activeCount int64
    err = db.Model(&models.EnrollmentTable{}).
        Where("student_id = ? AND semester = ? AND status = ?", studentID, maxSemester, models.StatusActive).
        Count(&activeCount).Error

    if err != nil {
        return 0, err
    }

    // If they still have active classes, they are currently in that semester
    if activeCount > 0 {
        return maxSemester, nil
    }

    // If all classes in that semester are finished (Passed or Failed), they advance to the next semester
    return maxSemester + 1, nil
}

func getPassedSubjects(db *gorm.DB, studentID uint) (map[uint]bool, error) {
	var passedIDs []uint
	err := db.Model(&models.EnrollmentTable{}).
		Joins("JOIN subject_groups ON subject_groups.id = enrollment_tables.subject_group_id").
		Where("enrollment_tables.student_id = ? AND enrollment_tables.status = ?", studentID, models.StatusPassed).
		Pluck("subject_groups.subject_id", &passedIDs).Error

	if err != nil {
		return nil, err
	}

	passedMap := make(map[uint]bool)
	for _, id := range passedIDs {
		passedMap[id] = true
	}
	return passedMap, nil
}

func evaluatePrerequisites(db *gorm.DB, choosable []models.ChoosableSubjects, passedMap map[uint]bool) ([]models.ChoosableSubjects, []uint, error) {
	var candidateIDs []uint
	for _, cs := range choosable {
		// Exclude subjects the student has already passed
		if !passedMap[cs.SubjectID] {
			candidateIDs = append(candidateIDs, cs.SubjectID)
		}
	}

	if len(candidateIDs) == 0 {
		return nil, nil, nil
	}

	var prereqs []models.SubjectPrerequisite
	if err := db.Where("subject_id IN ?", candidateIDs).Find(&prereqs).Error; err != nil {
		return nil, nil, err
	}

	// prereqMap[SubjectID] => map[BucketID][]PrerequisiteSubjectId
	prereqMap := make(map[uint]map[uint][]uint)
	for _, p := range prereqs {
		if prereqMap[p.SubjectId] == nil {
			prereqMap[p.SubjectId] = make(map[uint][]uint)
		}
		prereqMap[p.SubjectId][p.BucketId] = append(prereqMap[p.SubjectId][p.BucketId], p.SubjectPrerequisiteId)
	}

	var validChoosable []models.ChoosableSubjects
	var validIDs []uint

	for _, cs := range choosable {
		if passedMap[cs.SubjectID] {
			continue // Skip already passed subjects
		}

		buckets, hasPrereqs := prereqMap[cs.SubjectID]
		isUnlocked := !hasPrereqs

		if hasPrereqs {
			for _, requiredIDs := range buckets {
				bucketSatisfied := true
				for _, reqID := range requiredIDs {
					if !passedMap[reqID] {
						bucketSatisfied = false
						break
					}
				}
				if bucketSatisfied {
					isUnlocked = true
					break
				}
			}
		}

		if isUnlocked {
			validChoosable = append(validChoosable, cs)
			validIDs = append(validIDs, cs.SubjectID)
		}
	}

	return validChoosable, validIDs, nil
}

func buildSubjectResponseDTO(db *gorm.DB, choosable []models.ChoosableSubjects, subjectIDs []uint) ([]models.ChoosableSubjectResponse, error) {
	var allGroups []models.SubjectGroup
	if err := db.Where("subject_id IN ?", subjectIDs).Preload("Lecturer").Find(&allGroups).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch subject groups: %w", err)
	}

	groupsBySubject := make(map[uint][]models.SubjectGroup)
	for _, group := range allGroups {
		groupsBySubject[group.SubjectId] = append(groupsBySubject[group.SubjectId], group)
	}

	var response []models.ChoosableSubjectResponse
	for _, cs := range choosable {
		groups := groupsBySubject[cs.SubjectID]

		response = append(response, models.ChoosableSubjectResponse{
			SubjectData:   cs.Subject,
			SubjectGroups: groups,
		})
	}

	return response, nil
}