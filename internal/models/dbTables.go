package models

type Status int

const (
	StatusActive Status = iota + 1
	StatusPassed
	StatusFailed
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

type Semester int

const (
	SemesterFall Semester = iota + 1
	SemesterSpring
)

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

type EnrollmentTable struct {
	ID             int32    `gorm:"column:id;primaryKey;autoIncrement"`
	Status         Status   `gorm:"column:status;type:smallint;not null"`
	SubjectGroupID int32    `gorm:"column:subjectGroupId;type:integer;not null"`
	Year           int16    `gorm:"column:year;type:smallint;not null"`
	Semester       Semester `gorm:"column:semester;type:smallint;not null"`
	StudentID      int32    `gorm:"column:studentId;type:integer;not null"`
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
	ID          int32       `gorm:"column:id;primaryKey;autoIncrement"`
	FacultyID   int32       `gorm:"column:facultyId;type:integer;not null"`
	SubjectID   int32       `gorm:"column:subjectId;type:integer;not null"`
	SubjectType SubjectType `gorm:"column:subjectType;type:smallint;not null"`
}

type Faculty struct {
	ID          int32  `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `gorm:"column:name;type:varchar(255);not null"`
	Description string `gorm:"column:description;type:varchar(255)"`
}

type Subject struct {
	ID             int32   `gorm:"column:id;primaryKey;autoIncrement"`
	Name           string  `gorm:"column:name;type:varchar(255);not null"`
	Description    string  `gorm:"column:description;type:varchar(255)"`
	Credit         int32   `gorm:"column:credit;type:integer;not null"`
	FacultyID      int32   `gorm:"column:facultyId;type:integer;not null"`
	SyllabusPdfURL *string `gorm:"column:syllabusPdfURL;type:text"`
}

type SubjectPrerequisite struct {
	SubjectId           int16 `gorm:"column:subjectId;primaryKey"`
	SubjectPrerequisite int16 `gorm:"column:subjectPrerequisite;primaryKey"`
	BucketId            int16 `gorm:"column:bucketId"`
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
		return "cordinator"
	case RoleManager:
		return "manager"
	case RoleAdmin:
		return "admin"
	default:
		return ""
	}
}

type Users struct {
	Id            int16  `gorm:"column:id;primaryKey"`
	Email         string `gorm:"column:email"`
	Password      string `gorm:"column:password"`
	Role          Role   `gorm:"column:role"`
	FirstName     string `gorm:"column:firstName"`
	LastName      string `gorm:"column:lastName"`
	Gender        string `gorm:"column:gender"`
	PersonalNo    string `gorm:"column:personalNo"`
	PhoneNumber   string `gorm:"column:phoneNumber"`
	PhotoUrl      string `gorm:"column:photoURL"`
	AcademicPause *bool  `gorm:"column:academicPause"`
}

type Lecturer struct {
	Id          int16  `gorm:"column:id;primaryKey"`
	FirstName   string `gorm:"column:firstName"`
	LastName    string `gorm:"column:lastName"`
	Description string `gorm:"column:description"`
	PhotoUrl    string `gorm:"column:photoURL"`
}

type LectureType int

const (
	Lecture = iota + 1
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
	Id          int16       `gorm:"column:id;primaryKey"`
	GroupNo     int16       `gorm:"column:groupNo"`
	SubjectId   int16       `gorm:"column:subjectId"`
	Day         int8        `gorm:"column:day"`
	From        string      `gorm:"column:from"`
	To          string      `gorm:"column:to"`
	MaxStudent  int8        `gorm:"column:maxStudent"`
	LectureType LectureType `gorm:"column:lectureType"`
	LecturerId  int16       `gorm:"column:lecturerId"`
}

type StudentPriorityList struct {
	Id             int16 `gorm:"column:id;primaryKey"`
	Priority       int16 `gorm:"column:priority"`
	SubjectGroupId int   `gorm:"column:subjectGroupId"`
	StudentId      int   `gorm:"column:studentId"`
}
