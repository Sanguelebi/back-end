package models

type Status int
type Semester int

const (
	StatusActive Status = iota + 1
	StatusPassed
	StatusFailed
)

const (
	SemesterFall Semester = iota + 1
	SemesterSpring
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusPassed:
		return "passed"
	case StatusFailed:
		return "failed"
	default:
		return ""
	}
}

func (s Semester) String() string {
	switch s {
	case SemesterFall:
		return "fall"
	case SemesterSpring:
		return "spring"
	default:
		return ""
	}
}

type Faculty struct {
	ID          int32  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:varchar(255);not null"`
}

type Users struct {
	ID            int32  `gorm:"primaryKey;autoIncrement"`
	Email         string `gorm:"type:varchar(254);not null;unique"`
	Password      string `gorm:"type:varchar(255);not null"`
	Role          Role   `gorm:"type:smallint;not null"`
	FirstName     string `gorm:"type:varchar(255);not null"`
	LastName      string `gorm:"type:varchar(255);not null"`
	Gender        string `gorm:"type:varchar(255);not null"`
	PersonalNo    string `gorm:"type:varchar(255);not null"`
	PhoneNumber   string `gorm:"type:varchar(255);not null"`
	PhotoUrl      string `gorm:"type:text;not null"`
	AcademicPause *bool  `gorm:"type:boolean"`
}

type Subject struct {
	ID             int32   `gorm:"primaryKey;autoIncrement"`
	Name           string  `gorm:"type:varchar(255);not null"`
	Description    string  `gorm:"type:varchar(255);not null"`
	Credit         int32   `gorm:"type:integer;not null"`
	FacultyID      int32   `gorm:"type:integer;not null"`
	Faculty        Faculty `gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	SyllabusPdfURL *string `gorm:"type:text"`
}

type SubjectPrerequisite struct {
	SubjectId             int32   `gorm:"primaryKey;type:integer;not null"`
	Subject               Subject `gorm:"foreignKey:SubjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubjectPrerequisiteId int32   `gorm:"primaryKey;type:integer;not null"`
	PrerequisiteSubject   Subject `gorm:"foreignKey:SubjectPrerequisiteId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BucketId              int32   `gorm:"type:integer;not null"`
}

type EnrollmentTable struct {
	ID             int32        `gorm:"primaryKey;autoIncrement"`
	Status         Status       `gorm:"type:smallint;not null"`
	SubjectGroupID int32        `gorm:"type:integer;not null"`
	SubjectGroup   SubjectGroup `gorm:"foreignKey:SubjectGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Year           int16        `gorm:"type:smallint;not null"`
	Semester       Semester     `gorm:"type:smallint;not null"`
	StudentID      int32        `gorm:"type:integer;not null"`
	Student        Users        `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type SubjectType int

const (
	SubjectProfessional SubjectType = iota + 1
	SubjectFree
)

func (s SubjectType) String() string {
	switch s {
	case SubjectProfessional:
		return "professional"
	case SubjectFree:
		return "free"
	default:
		return ""
	}
}

type ChoosableSubjects struct {
	ID          int32       `gorm:"primaryKey;autoIncrement"`
	FacultyID   int32       `gorm:"type:integer;not null"`
	Faculty     Faculty     `gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubjectID   int32       `gorm:"type:integer;not null"`
	Subject     Subject     `gorm:"foreignKey:SubjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubjectType SubjectType `gorm:"type:smallint;not null"`
}

type Role int

const (
	RoleStudent Role = iota + 1
	RoleCordinator
	RoleManager
	RoleAdmin
)

func (s Role) String() string {
	switch s {
	case RoleStudent:
		return "student"
	case RoleCordinator:
		return "coordinator"
	case RoleManager:
		return "manager"
	case RoleAdmin:
		return "admin"
	default:
		return ""
	}
}

type Lecturer struct {
	ID          int32   `gorm:"primaryKey;autoIncrement"`
	FirstName   string  `gorm:"type:varchar(255);not null"`
	LastName    string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text;not null"`
	PhotoUrl    *string `gorm:"type:text"`
}

type LectureType int

const (
	Lecture LectureType = iota + 1
	Seminar
)

func (s LectureType) String() string {
	switch s {
	case Lecture:
		return "lecture"
	case Seminar:
		return "seminar"
	default:
		return ""
	}
}

type SubjectGroup struct {
	ID          int32       `gorm:"primaryKey;autoIncrement"`
	GroupNo     int16       `gorm:"type:smallint;not null"`
	SubjectId   int32       `gorm:"type:integer;not null"`
	Subject     Subject     `gorm:"foreignKey:SubjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Day         int8        `gorm:"type:smallint;not null"`
	From        string      `gorm:"type:varchar(5);not null"`
	To          string      `gorm:"type:varchar(5);not null"`
	MaxStudent  int8        `gorm:"type:smallint;not null"`
	LectureType LectureType `gorm:"type:smallint;not null"`
	LecturerId  int32       `gorm:"type:integer;not null"`
	Lecturer    Lecturer    `gorm:"foreignKey:LecturerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type StudentPriorityList struct {
	ID             int32        `gorm:"primaryKey;autoIncrement"`
	Priority       int16        `gorm:"type:smallint;not null"`
	SubjectGroupId int32        `gorm:"type:integer;not null"`
	SubjectGroup   SubjectGroup `gorm:"foreignKey:SubjectGroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	StudentId      int32        `gorm:"type:integer;not null"`
	Student        Users        `gorm:"foreignKey:StudentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type LecturerRating struct {
	ID         	uint 		`gorm:"primaryKey;autoIncrement"`
	LecturerId 	uint 		`gorm:"type:integer;not null"`
	Lecturer	Lecturer	`gorm:"foreignKey:LecturerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rating     	int16 		`gorm:"type:integer;not null"`
}
