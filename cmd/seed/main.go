package main

import (
	"log"
	"os"
	"time"

	"course-selection-platform/config"
	"course-selection-platform/database"
	"course-selection-platform/internal/models"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env values first
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return
	}

	// Load the DB Config first and then connect to the DB with it.
	cfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Config initialization has failed: %v", err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	log.Println("Connected to the DB Successfully. Wiping existing data...")

	// ======== 0. Wipe All Existing Data ========
	// Truncate tables in reverse order of foreign key dependencies and reset auto-increment IDs.
	wipeSQL := `
		TRUNCATE TABLE 
			active_registrations,
			inactive_registrations,
			subject_difficulties,
			subject_recommendations,
			lecturer_ratings,
			enrollment_tables,
			student_priority_lists,
			subject_groups,
			lecturers,
			choosable_subjects,
			subject_prerequisites,
			subjects,
			users,
			faculties
		RESTART IDENTITY CASCADE;
	`
	if err := db.Exec(wipeSQL).Error; err != nil {
		log.Fatalf("Failed to wipe database tables: %v", err)
	}
	log.Println("Database wiped clean. Seeding all data...")

	// ======== 1. Seed Faculties ======== 
	faculties := []models.Faculty{
		{Name: "კომპიუტერული მეცნიერება და IT", Description: "კომპიუტერული მეცნიერებისა და ინფორმაციული ტექნოლოგიების ფაკულტეტი"},
		{Name: "ბიზნესის ადმინისტრირება", Description: "ეკონომიკისა და ბიზნესის ფაკულტეტი"},
	}
	for i := range faculties {
		db.FirstOrCreate(&faculties[i], models.Faculty{Name: faculties[i].Name})
	}

	// ======== 2. Seed Users with Dynamic Bcrypt Hashing ======== 
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
			log.Fatalf("Failed to hash password for user %s: %v", seed.user.Email, err)
		}

		u := seed.user
		u.Password = string(hashedBytes)

		if err := db.FirstOrCreate(&u, models.Users{Email: u.Email}).Error; err != nil {
			log.Fatalf("Failed to seed user %s: %v", u.Email, err)
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
	}
	for i := range subjects {
		db.FirstOrCreate(&subjects[i], models.Subject{Name: subjects[i].Name})
	}

	// ======== 4. Seed Subject Prerequisites ========
	prerequisites := []models.SubjectPrerequisite{
		{
			SubjectId:             subjects[1].ID, // Database Systems requires...
			SubjectPrerequisiteId: subjects[0].ID, // Data Structures and Algorithms
			BucketId:              1,
		},
	}
	for i := range prerequisites {
		db.FirstOrCreate(&prerequisites[i], models.SubjectPrerequisite{
			SubjectId:             prerequisites[i].SubjectId,
			SubjectPrerequisiteId: prerequisites[i].SubjectPrerequisiteId,
		})
	}

	// ======== 5. Seed Choosable Subjects ========
	choosableSubjects := []models.ChoosableSubjects{
		{
			FacultyID:   faculties[0].ID,
			SubjectID:   subjects[0].ID,
			SubjectType: models.SubjectProfessional,
		},
		{
			FacultyID:   faculties[0].ID,
			SubjectID:   subjects[1].ID,
			SubjectType: models.SubjectProfessional,
		},
	}
	for i := range choosableSubjects {
		db.FirstOrCreate(&choosableSubjects[i], models.ChoosableSubjects{
			FacultyID: choosableSubjects[i].FacultyID,
			SubjectID: choosableSubjects[i].SubjectID,
		})
	}

	// ======== 6. Seed Lecturers ======== 
	lecturers := []models.Lecturer{
		{
			FirstName:   "ლევან",
			LastName:    "ქავთარაძე",
			Description: "პროგრამული უზრუნველყოფის ინჟინერიის პროფესორი",
		},
	}
	for i := range lecturers {
		db.FirstOrCreate(&lecturers[i], models.Lecturer{
			FirstName: lecturers[i].FirstName,
			LastName:  lecturers[i].LastName,
		})
	}

	// ======== 7. Seed Subject Groups ======== 
	subjectGroups := []models.SubjectGroup{
		{
			GroupNo:     1,
			SubjectId:   subjects[0].ID,
			Day:         1, // ორშაბათი
			From:        "10:00",
			To:          "12:00",
			MaxStudent:  30,
			LectureType: models.Lecture,
			LecturerId:  lecturers[0].ID,
		},
	}
	for i := range subjectGroups {
		db.FirstOrCreate(&subjectGroups[i], models.SubjectGroup{
			SubjectId: subjectGroups[i].SubjectId,
			GroupNo:   subjectGroups[i].GroupNo,
		})
	}

	// ======== 8. Seed Student Priority List ======== 
	priorities := []models.StudentPriorityList{
		{
			Priority:       1,
			SubjectGroupId: subjectGroups[0].ID,
			StudentId:      users[1].ID,
		},
	}
	for i := range priorities {
		db.FirstOrCreate(&priorities[i], models.StudentPriorityList{
			StudentId:      priorities[i].StudentId,
			SubjectGroupId: priorities[i].SubjectGroupId,
		})
	}

	// ======== 9. Seed Enrollments ======== 
	enrollments := []models.EnrollmentTable{
		{
			Status:         models.StatusActive,
			SubjectGroupID: subjectGroups[0].ID,
			Year:           2026,
			Semester:       models.SemesterFall,
			StudentID:      users[1].ID, // student user
		},
	}
	for i := range enrollments {
		db.FirstOrCreate(&enrollments[i], models.EnrollmentTable{
			StudentID:      enrollments[i].StudentID,
			SubjectGroupID: enrollments[i].SubjectGroupID,
		})
	}

	// ======== 10. Seed Lecturer Ratings ======== 
	ratings := []models.LecturerRating{
		{
			LecturerId: lecturers[0].ID,
			Rating:     5,
		},
	}
	for i := range ratings {
		db.FirstOrCreate(&ratings[i], models.LecturerRating{
			LecturerId: ratings[i].LecturerId,
		})
	}

	// ======== 11. Seed Subject Recommendations (Votes) ======== 
	recommendations := []models.SubjectRecommendation{
		{
			SubjectId: subjects[0].ID,
			UserId:    users[1].ID,
			VoteType:  models.Like,
		},
	}
	for i := range recommendations {
		db.FirstOrCreate(&recommendations[i], models.SubjectRecommendation{
			SubjectId: recommendations[i].SubjectId,
			UserId:    recommendations[i].UserId,
		})
	}

	// ======== 12. Seed Subject Difficulties ======== 
	difficulties := []models.SubjectDifficulty{
		{
			SubjectId:       subjects[0].ID,
			DifficultyScale: 4,
		},
		{
			SubjectId:       subjects[1].ID,
			DifficultyScale: 3,
		},
	}
	for i := range difficulties {
		db.FirstOrCreate(&difficulties[i], models.SubjectDifficulty{
			SubjectId: difficulties[i].SubjectId,
		})
	}

	// ======== 13. Seed Registration Windows (Active & Inactive) ======== 
	activeReg := models.ActiveRegistration{
		From: time.Now(),
		To:   time.Now().Add(14 * 24 * time.Hour),
	}
	db.Create(&activeReg)

	inactiveReg := models.InactiveRegistration{
		From: time.Now().Add(-30 * 24 * time.Hour),
		To:   time.Now().Add(-15 * 24 * time.Hour),
	}
	db.Create(&inactiveReg)

	log.Println("მონაცემთა ბაზის წაშლა და სიდინგი (Seed) წარმატებით დასრულდა!")
}