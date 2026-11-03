package models

import "time"

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
	ID          uint  	`gorm:"primaryKey;autoIncrement"`
	Name        string 	`gorm:"type:varchar(255);not null"`
	Description string 	`gorm:"type:varchar(255);not null"`
}

type StudentStatus int

const (
	Active StudentStatus = iota + 1
	Inactive
	AcademicPause
)

func (s StudentStatus) String() string {
	switch s {
	case Active: return "active"
	case Inactive: return "inactive"
	case AcademicPause: return "academicPause"
	default: return ""
	}
}

type Gender int

const (
	Male Gender = iota + 1
	Female
)

type Users struct {
	ID            uint  			`gorm:"primaryKey;autoIncrement"`
	Email         string 			`gorm:"type:varchar(254);not null;unique"`
	Password      string 			`gorm:"type:varchar(255);not null"`
	Role          Role   			`gorm:"type:smallint;not null"`
	FirstName     string 			`gorm:"type:varchar(255);not null"`
	LastName      string 			`gorm:"type:varchar(255);not null"`
	Gender        Gender 			`gorm:"type:smallint;not null"`
	PersonalNo    string 			`gorm:"type:varchar(255);not null"`
	PhoneNumber   string 			`gorm:"type:varchar(255);not null"`
	PhotoUrl      string 			`gorm:"type:text;not null"`
	StudentStatus *StudentStatus  	`gorm:"type:smallint"`
}

type Subject struct {
	ID             uint   	`gorm:"primaryKey;autoIncrement"`
	Name           string  	`gorm:"type:varchar(255);not null"`
	Description    string  	`gorm:"type:varchar(255);not null"`
	Credit         int32   	`gorm:"type:integer;not null"`
	FacultyID      uint   	`gorm:"type:integer;not null"`
	Faculty        Faculty 	`gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	SyllabusPdfURL *string 	`gorm:"type:text"`
}

type SubjectPrerequisite struct {
	SubjectId             uint   	`gorm:"primaryKey;type:integer;not null"`
	Subject               Subject 	`gorm:"foreignKey:SubjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubjectPrerequisiteId uint   	`gorm:"primaryKey;type:integer;not null"`
	PrerequisiteSubject   Subject 	`gorm:"foreignKey:SubjectPrerequisiteId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BucketId              uint   	`gorm:"type:integer;not null"`
}

type EnrollmentTable struct {
	ID             uint        	`gorm:"primaryKey;autoIncrement"`
	Status         Status       `gorm:"type:smallint;not null"`
	SubjectGroupID uint        	`gorm:"type:integer;not null"`
	SubjectGroup   SubjectGroup `gorm:"foreignKey:SubjectGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Year           int16        `gorm:"type:smallint;not null"`
	Semester       Semester     `gorm:"type:smallint;not null"`
	StudentID      uint        	`gorm:"type:integer;not null"`
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
	ID          uint       	`gorm:"primaryKey;autoIncrement"`
	FacultyID   uint       	`gorm:"type:integer;not null"`
	Faculty     Faculty     `gorm:"foreignKey:FacultyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SubjectID   uint       	`gorm:"type:integer;not null"`
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
	ID          uint   `gorm:"primaryKey;autoIncrement"`
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
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	GroupNo     int16       `gorm:"type:smallint;not null"`
	SubjectId   uint       `gorm:"type:integer;not null"`
	Subject     Subject     `gorm:"foreignKey:SubjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Day         int8        `gorm:"type:smallint;not null"`
	From        string      `gorm:"type:varchar(5);not null"`
	To          string      `gorm:"type:varchar(5);not null"`
	MaxStudent  int8        `gorm:"type:smallint;not null"`
	LectureType LectureType `gorm:"type:smallint;not null"`
	LecturerId  uint       `gorm:"type:integer;not null"`
	Lecturer    Lecturer    `gorm:"foreignKey:LecturerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

type StudentPriorityList struct {
	ID             uint        	`gorm:"primaryKey;autoIncrement"`
	Priority       int16        `gorm:"type:smallint;not null"`
	SubjectGroupId uint        	`gorm:"type:integer;not null"`
	SubjectGroup   SubjectGroup `gorm:"foreignKey:SubjectGroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	StudentId      uint        	`gorm:"type:integer;not null"`
	Student        Users        `gorm:"foreignKey:StudentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type LecturerRating struct {
	ID         	uint 		`gorm:"primaryKey;autoIncrement"`
	LecturerId 	uint 		`gorm:"type:integer;not null"`
	Lecturer	Lecturer	`gorm:"foreignKey:LecturerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rating     	int16 		`gorm:"type:integer;not null"`
}

type Vote int

const (
	Like Vote = iota + 1
	Dislike
)

func (s Vote) String() string {
	switch s {
	case Like:
		return "like"
	case Dislike:
		return "dislike"
	default:
		return ""
	}
}

type SubjectRecommendation struct {
	ID        uint   	`gorm:"primaryKey;autoIncrement"`
	SubjectId uint   	`gorm:"type:integer;not null"`
	Subject   Subject 	`gorm:"foreignKey:SubjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserId    uint   	`gorm:"type:integer;not null"`
	VoteType  Vote    	`gorm:"not null"`
}

type SubjectDifficulty struct {
	ID              uint   	`gorm:"primaryKey;autoIncrement"`
	SubjectId       uint   	`gorm:"type:integer;not null"`
	Subject         Subject `gorm:"foreignKey:SubjectId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DifficultyScale int16   `gorm:"type:integer;not null"`
}

type ActiveRegistration struct {
	ID		uint 		`gorm:"primaryKey;autoIncrement"`
	From 	time.Time
	To		time.Time
}

type InactiveRegistration struct {
	ID		uint 		`gorm:"primaryKey;autoIncrement"`
	From 	time.Time
	To		time.Time
}