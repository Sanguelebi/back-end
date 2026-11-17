package database

import (
	"course-selection-platform/config"
	"course-selection-platform/internal/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB will accept the DBConfig as a dependency.
func ConnectDB(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying database connection: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.AutoMigrate(
		&models.Faculty{},
		&models.Users{},
		&models.Subject{},
		&models.EnrollmentTable{},
		&models.ChoosableSubjects{},
		&models.SubjectPrerequisite{},
		&models.Lecturer{},
		&models.SubjectGroup{},
		&models.StudentPriorityList{},
		&models.LecturerRating{},
		&models.SubjectDifficulty{},
		&models.SubjectRecommendation{},
		&models.ActiveRegistration{},
		&models.InactiveRegistration{},
		&models.StudentSubjectAllowance{},
	)

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)

	return db, nil
}
