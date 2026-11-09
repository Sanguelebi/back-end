package database

import (
	"log"
	"time"

	"course-selection-platform/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase checks if the database is empty. If no users exist,
// it wipes dependent tables and populates the initial seed dataset inside a transaction.
func SeedDatabase(db *gorm.DB) error {
	var count int64

	// Check if the database already contains seeded data
	if err := db.Model(&models.Users{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("[INFO] Database already contains data. Skipping seed execution.")
		return nil
	}

	log.Println("[INFO] Database is empty. Starting transactional database seeding...")

	// Execute all wipe and seed operations inside a DB transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// ======== 0. Wipe All Existing Data ========
		tables, err := tx.Migrator().GetTables()
		if err != nil {
			return err
		}

		for _, table := range tables {
			if err := tx.Migrator().DropTable(table); err != nil {
				return err
			}
		}

		if err := tx.AutoMigrate(
			&models.Faculty{},
			&models.Users{},
			&models.Subject{},
			&models.SubjectPrerequisite{},
			&models.ChoosableSubjects{},
			&models.Lecturer{},
			&models.SubjectGroup{},
			&models.StudentPriorityList{},
			&models.EnrollmentTable{},
			&models.LecturerRating{},
			&models.SubjectRecommendation{},
			&models.SubjectDifficulty{},
			&models.ActiveRegistration{},
			&models.InactiveRegistration{},
		); err != nil {
			return err
		}

		log.Println("[INFO] Database dropped and re-migrated cleanly.")

		// ======== 1. Seed Faculties ========
		faculties := []models.Faculty{
			{Name: "კომპიუტერული მეცნიერება და IT", Description: "კომპიუტერული მეცნიერებისა და ინფორმაციული ტექნოლოგიების ფაკულტეტი"},
			{Name: "ბიზნესის ადმინისტრირება", Description: "ეკონომიკისა და ბიზნესის ფაკულტეტი"},
		}
		for i := range faculties {
			if err := tx.FirstOrCreate(&faculties[i], models.Faculty{Name: faculties[i].Name}).Error; err != nil {
				return err
			}
		}

		// ======== 2. Seed Users ========
		activeStudentStatus := models.Active

		userSeeds := []struct {
			user          models.Users
			plainPassword string
		}{
			{
				user: models.Users{
					Email:         "admin@university.edu",
					Role:          models.RoleAdmin,
					FirstName:     "სისტემური",
					LastName:      "ადმინისტრატორი",
					Gender:        models.Male,
					PersonalNo:    "01001001001",
					PhoneNumber:   "+995555123456",
					PhotoUrl:      "https://example.com/admin.jpg",
					StudentStatus: nil,
				},
				plainPassword: "AdminSecret123!",
			},
			{
				user: models.Users{
					Email:         "student@university.edu",
					Role:          models.RoleStudent,
					FirstName:     "გიორგი",
					LastName:      "ბერიძე",
					Gender:        models.Male,
					PersonalNo:    "01001001002",
					PhoneNumber:   "+995555654321",
					PhotoUrl:      "https://example.com/student.jpg",
					StudentStatus: &activeStudentStatus,
				},
				plainPassword: "StudentPass123!",
			},
		}

		var users []models.Users
		for _, seed := range userSeeds {
			hashedBytes, err := bcrypt.GenerateFromPassword([]byte(seed.plainPassword), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			u := seed.user
			u.Password = string(hashedBytes)

			if err := tx.FirstOrCreate(&u, models.Users{Email: u.Email}).Error; err != nil {
				return err
			}
			users = append(users, u)
		}

		// ======== 3. Seed Subjects ========
		subjects := []models.Subject{
			{
				Name:        "მონაცემთა სტრუქტურები და ალგორითმები",
				Description: "მოწინავე ალგორითმები და მონაცემთა სტრუქტურები Go-ში",
				Credit:      6,
				FacultyID:   faculties[0].ID,
			},
			{
				Name:        "მონაცემთა ბაზების სისტემები",
				Description: "რელაციური მონაცემთა ბაზები, SQL და GORM ORM",
				Credit:      5,
				FacultyID:   faculties[0].ID,
			},
			{
				Name:        "ვებ-პროგრამირება (Go & React)",
				Description: "სრული დასტის (Full-Stack) ვებ აპლიკაციების შემუშავება",
				Credit:      6,
				FacultyID:   faculties[0].ID,
			},
			{
				Name:        "კომპიუტერული ქსელები",
				Description: "TCP/IP პროტოკოლები, როუტინგი და ქსელური უსაფრთხოება",
				Credit:      5,
				FacultyID:   faculties[0].ID,
			},
			{
				Name:        "მიკროეკონომიკის პრინციპები",
				Description: "საბაზრო მექანიზმები, მოთხოვნა-მიწოდება და ფასწარმოქმნა",
				Credit:      5,
				FacultyID:   faculties[1].ID,
			},
			{
				Name:        "ფინანსური აღრიცხვა",
				Description: "ბუღალტრული აღრიცხვის საფუძვლები და ფინანსური უწყისები",
				Credit:      6,
				FacultyID:   faculties[1].ID,
			},
		}
		for i := range subjects {
			if err := tx.FirstOrCreate(&subjects[i], models.Subject{Name: subjects[i].Name}).Error; err != nil {
				return err
			}
		}

		// ======== 4. Seed Subject Prerequisites ========
		prerequisites := []models.SubjectPrerequisite{
			{
				SubjectId:             subjects[1].ID,
				SubjectPrerequisiteId: subjects[0].ID,
				BucketId:              1,
			},
			{
				SubjectId:             subjects[2].ID,
				SubjectPrerequisiteId: subjects[1].ID,
				BucketId:              1,
			},
		}
		for i := range prerequisites {
			if err := tx.FirstOrCreate(&prerequisites[i], models.SubjectPrerequisite{
				SubjectId:             prerequisites[i].SubjectId,
				SubjectPrerequisiteId: prerequisites[i].SubjectPrerequisiteId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 5. Seed Choosable Subjects ========
		choosableSubjects := []models.ChoosableSubjects{
			{FacultyID: faculties[0].ID, SubjectID: subjects[0].ID, SubjectType: models.SubjectProfessional},
			{FacultyID: faculties[0].ID, SubjectID: subjects[1].ID, SubjectType: models.SubjectProfessional},
			{FacultyID: faculties[0].ID, SubjectID: subjects[2].ID, SubjectType: models.SubjectProfessional},
			{FacultyID: faculties[0].ID, SubjectID: subjects[3].ID, SubjectType: models.SubjectFree},
			{FacultyID: faculties[1].ID, SubjectID: subjects[4].ID, SubjectType: models.SubjectProfessional},
			{FacultyID: faculties[1].ID, SubjectID: subjects[5].ID, SubjectType: models.SubjectProfessional},
		}
		for i := range choosableSubjects {
			if err := tx.FirstOrCreate(&choosableSubjects[i], models.ChoosableSubjects{
				FacultyID: choosableSubjects[i].FacultyID,
				SubjectID: choosableSubjects[i].SubjectID,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 6. Seed Lecturers ========
		lecturers := []models.Lecturer{
			{
				FirstName:   "ლევან",
				LastName:    "ქავთარაძე",
				Description: "პროგრამული უზრუნველყოფის ინჟინერიის პროფესორი",
			},
			{
				FirstName:   "ნინო",
				LastName:    "დოლიძე",
				Description: "მონაცემთა ბაზებისა და ვებ-ტექნოლოგიების ლექტორი",
			},
			{
				FirstName:   "გიორგი",
				LastName:    "მაისურაძე",
				Description: "ბიზნესის ადმინისტრირებისა და ეკონომიკის დოქტორი",
			},
		}
		for i := range lecturers {
			if err := tx.FirstOrCreate(&lecturers[i], models.Lecturer{
				FirstName: lecturers[i].FirstName,
				LastName:  lecturers[i].LastName,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 7. Seed Subject Groups ========
		subjectGroups := []models.SubjectGroup{
			{
				GroupNo:     1,
				SubjectId:   subjects[0].ID,
				Day:         1,
				From:        "10:00",
				To:          "12:00",
				MaxStudent:  30,
				LectureType: models.Lecture,
				LecturerId:  lecturers[0].ID,
			},
			{
				GroupNo:     2,
				SubjectId:   subjects[0].ID,
				Day:         2,
				From:        "12:00",
				To:          "14:00",
				MaxStudent:  25,
				LectureType: models.Seminar,
				LecturerId:  lecturers[0].ID,
			},
			{
				GroupNo:     1,
				SubjectId:   subjects[1].ID,
				Day:         3,
				From:        "14:00",
				To:          "16:00",
				MaxStudent:  30,
				LectureType: models.Lecture,
				LecturerId:  lecturers[1].ID,
			},
			{
				GroupNo:     2,
				SubjectId:   subjects[1].ID,
				Day:         4,
				From:        "10:00",
				To:          "12:00",
				MaxStudent:  20,
				LectureType: models.Seminar,
				LecturerId:  lecturers[1].ID,
			},
			{
				GroupNo:     1,
				SubjectId:   subjects[2].ID,
				Day:         5,
				From:        "16:00",
				To:          "18:00",
				MaxStudent:  25,
				LectureType: models.Lecture,
				LecturerId:  lecturers[1].ID,
			},
			{
				GroupNo:     1,
				SubjectId:   subjects[3].ID,
				Day:         1,
				From:        "14:00",
				To:          "16:00",
				MaxStudent:  30,
				LectureType: models.Lecture,
				LecturerId:  lecturers[0].ID,
			},
			{
				GroupNo:     1,
				SubjectId:   subjects[4].ID,
				Day:         2,
				From:        "10:00",
				To:          "12:00",
				MaxStudent:  40,
				LectureType: models.Lecture,
				LecturerId:  lecturers[2].ID,
			},
			{
				GroupNo:     1,
				SubjectId:   subjects[5].ID,
				Day:         4,
				From:        "14:00",
				To:          "16:00",
				MaxStudent:  35,
				LectureType: models.Lecture,
				LecturerId:  lecturers[2].ID,
			},
		}
		for i := range subjectGroups {
			if err := tx.FirstOrCreate(&subjectGroups[i], models.SubjectGroup{
				SubjectId: subjectGroups[i].SubjectId,
				GroupNo:   subjectGroups[i].GroupNo,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 8. Seed Initial Student Priority List ========
		priorities := []models.StudentPriorityList{
			{
				Priority:       1,
				SubjectGroupId: subjectGroups[0].ID,
				StudentId:      users[1].ID,
			},
		}
		for i := range priorities {
			if err := tx.FirstOrCreate(&priorities[i], models.StudentPriorityList{
				StudentId:      priorities[i].StudentId,
				SubjectGroupId: priorities[i].SubjectGroupId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 9. Seed Enrollments ========
		enrollments := []models.EnrollmentTable{
			{
				Status:         models.StatusActive,
				SubjectGroupID: subjectGroups[0].ID,
				Year:           2026,
				Semester:       models.SemesterFall,
				StudentID:      users[1].ID,
			},
		}
		for i := range enrollments {
			if err := tx.FirstOrCreate(&enrollments[i], models.EnrollmentTable{
				StudentID:      enrollments[i].StudentID,
				SubjectGroupID: enrollments[i].SubjectGroupID,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 10. Seed Lecturer Ratings ========
		ratings := []models.LecturerRating{
			{LecturerId: lecturers[0].ID, Rating: 5},
			{LecturerId: lecturers[1].ID, Rating: 4},
			{LecturerId: lecturers[2].ID, Rating: 5},
		}
		for i := range ratings {
			if err := tx.FirstOrCreate(&ratings[i], models.LecturerRating{
				LecturerId: ratings[i].LecturerId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 11. Seed Subject Recommendations ========
		recommendations := []models.SubjectRecommendation{
			{SubjectId: subjects[0].ID, UserId: users[1].ID, VoteType: models.Like},
			{SubjectId: subjects[1].ID, UserId: users[1].ID, VoteType: models.Like},
		}
		for i := range recommendations {
			if err := tx.FirstOrCreate(&recommendations[i], models.SubjectRecommendation{
				SubjectId: recommendations[i].SubjectId,
				UserId:    recommendations[i].UserId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 12. Seed Subject Difficulties ========
		difficulties := []models.SubjectDifficulty{
			{SubjectId: subjects[0].ID, DifficultyScale: 4},
			{SubjectId: subjects[1].ID, DifficultyScale: 3},
			{SubjectId: subjects[2].ID, DifficultyScale: 4},
			{SubjectId: subjects[3].ID, DifficultyScale: 3},
		}
		for i := range difficulties {
			if err := tx.FirstOrCreate(&difficulties[i], models.SubjectDifficulty{
				SubjectId: difficulties[i].SubjectId,
			}).Error; err != nil {
				return err
			}
		}

		// ======== 13. Seed Registration Windows ========
		activeReg := models.ActiveRegistration{
			From: time.Now(),
			To:   time.Now().Add(14 * 24 * time.Hour),
		}
		if err := tx.Create(&activeReg).Error; err != nil {
			return err
		}

		inactiveReg := models.InactiveRegistration{
			From: time.Now().Add(-30 * 24 * time.Hour),
			To:   time.Now().Add(-15 * 24 * time.Hour),
		}
		if err := tx.Create(&inactiveReg).Error; err != nil {
			return err
		}

		log.Println("[INFO] Database seeded successfully!")
		return nil
	})
}
